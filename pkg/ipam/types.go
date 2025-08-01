// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package ipam

import (
	"log/slog"
	"net"

	"github.com/davecgh/go-spew/spew"

	agentK8s "github.com/cilium/cilium/daemon/k8s"
	"github.com/cilium/cilium/pkg/datapath/linux/sysctl"
	"github.com/cilium/cilium/pkg/datapath/types"
	"github.com/cilium/cilium/pkg/endpoint"
	"github.com/cilium/cilium/pkg/ipmasq"
	"github.com/cilium/cilium/pkg/k8s/client"
	"github.com/cilium/cilium/pkg/lock"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/node"
	"github.com/cilium/cilium/pkg/option"
)

// AllocationResult is the result of an allocation
type AllocationResult struct {
	// IP is the allocated IP
	IP net.IP

	// IPPoolName is the IPAM pool from which the above IP was allocated from
	IPPoolName Pool

	// CIDRs is a list of all CIDRs to which the IP has direct access to.
	// This is primarily useful if the IP has been allocated out of a VPC
	// subnet range and the VPC provides routing to a set of CIDRs in which
	// the IP is routable.
	CIDRs []string

	// PrimaryMAC is the MAC address of the primary interface. This is useful
	// when the IP is a secondary address of an interface which is
	// represented on the node as a Linux device and all routing of the IP
	// must occur through that master interface.
	PrimaryMAC string

	// GatewayIP is the IP of the gateway which must be used for this IP.
	// If the allocated IP is derived from a VPC, then the gateway
	// represented the gateway of the VPC or VPC subnet.
	GatewayIP string

	// ExpirationUUID is the UUID of the expiration timer. This field is
	// only set if AllocateNextWithExpiration is used.
	ExpirationUUID string

	// InterfaceNumber is a field for generically identifying an interface.
	// This is only useful in ENI mode.
	InterfaceNumber string
}

// Allocator is the interface for an IP allocator implementation
type Allocator interface {
	// Allocate allocates a specific IP or fails
	Allocate(ip net.IP, owner string, pool Pool) (*AllocationResult, error)

	// AllocateWithoutSyncUpstream allocates a specific IP without syncing
	// upstream or fails
	AllocateWithoutSyncUpstream(ip net.IP, owner string, pool Pool) (*AllocationResult, error)

	// Release releases a previously allocated IP or fails
	Release(ip net.IP, pool Pool) error

	// AllocateNext allocates the next available IP or fails if no more IPs
	// are available
	AllocateNext(owner string, pool Pool) (*AllocationResult, error)

	// AllocateNextWithoutSyncUpstream allocates the next available IP without syncing
	// upstream or fails if no more IPs are available
	AllocateNextWithoutSyncUpstream(owner string, pool Pool) (*AllocationResult, error)

	// Dump returns a map of all allocated IPs per pool with the IP represented as key in the
	// map. Dump must also provide a status one-liner to represent the overall status, e.g.
	// number of IPs allocated and overall health information if available.
	Dump() (map[Pool]map[string]string, string)

	// Capacity returns the total IPAM allocator capacity (not the current
	// available).
	Capacity() uint64

	// RestoreFinished marks the status of restoration as done
	RestoreFinished()
}

// IPAM is the configuration used for a particular IPAM type.
type IPAM struct {
	logger *slog.Logger

	nodeAddressing types.NodeAddressing
	config         *option.DaemonConfig

	IPv6Allocator Allocator
	IPv4Allocator Allocator

	// metadata provides information about a particular IP owner.
	metadata Metadata

	// owner maps an IP to the owner per pool.
	owner map[Pool]map[string]string

	// expirationTimers is a map of all expiration timers. Each entry
	// represents a IP allocation which is protected by an expiration
	// timer.
	expirationTimers map[timerKey]expirationTimer

	// mutex covers access to all members of this struct
	allocatorMutex lock.RWMutex

	// excludedIPS contains excluded IPs and their respective owners per pool. The key is a
	// combination pool:ip to avoid having to maintain a map of maps.
	excludedIPs map[string]string

	localNodeStore *node.LocalNodeStore
	k8sEventReg    K8sEventRegister
	nodeResource   agentK8s.LocalCiliumNodeResource
	mtuConfig      MtuConfiguration
	clientset      client.Clientset
	nodeDiscovery  Owner
	sysctl         sysctl.Sysctl
	ipMasqAgent    *ipmasq.IPMasqAgent
}

func (ipam *IPAM) EndpointCreated(ep *endpoint.Endpoint) {}

func (ipam *IPAM) EndpointDeleted(ep *endpoint.Endpoint, conf endpoint.DeleteConfig) {
	if !conf.NoIPRelease {
		if option.Config.EnableIPv4 {
			if err := ipam.ReleaseIP(ep.IPv4.AsSlice(), PoolOrDefault(ep.IPv4IPAMPool)); err != nil {
				ipam.logger.Warn("Unable to release IPv4 address during endpoint deletion", logfields.Error, err)
			}
		}
		if option.Config.EnableIPv6 {
			if err := ipam.ReleaseIP(ep.IPv6.AsSlice(), PoolOrDefault(ep.IPv6IPAMPool)); err != nil {
				ipam.logger.Warn("Unable to release IPv6 address during endpoint deletion", logfields.Error, err)
			}
		}
	}
}

func (ipam *IPAM) EndpointRestored(ep *endpoint.Endpoint) {}

// DebugStatus implements debug.StatusObject to provide debug status collection
// ability
func (ipam *IPAM) DebugStatus() string {
	ipam.allocatorMutex.RLock()
	str := spew.Sdump(
		"owners", ipam.owner,
		"expiration timers", ipam.expirationTimers,
		"excluded ips", ipam.excludedIPs,
	)
	ipam.allocatorMutex.RUnlock()
	return str
}

// Pool is the IP pool from which to allocate.
type Pool string

func (p Pool) String() string {
	return string(p)
}

type timerKey struct {
	ip   string
	pool Pool
}

type expirationTimer struct {
	uuid string
	stop chan<- struct{}
}

// LimitsNotFound is an error that signals lack of limits for given instance type
type LimitsNotFound struct{}

// Error implements error interface
func (_ LimitsNotFound) Error() string {
	return "Limits not found"
}
