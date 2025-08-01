// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package k8s

import (
	"fmt"
	"log/slog"
	"maps"
	"net"
	"net/netip"
	"slices"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"

	serviceStore "github.com/cilium/cilium/pkg/clustermesh/store"
	cmtypes "github.com/cilium/cilium/pkg/clustermesh/types"
	slim_corev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/core/v1"
	slim_discovery_v1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/discovery/v1"
	slim_metav1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	"github.com/cilium/cilium/pkg/k8s/types"
	"github.com/cilium/cilium/pkg/loadbalancer"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/metrics"
)

// EndpointSliceID identifies a Kubernetes EndpointSlice as well as the legacy
// v1.Endpoints.
type EndpointSliceID struct {
	ServiceName       loadbalancer.ServiceName
	EndpointSliceName string
}

// Endpoints is an abstraction for the Kubernetes endpoints object. Endpoints
// consists of a set of backend IPs in combination with a set of ports and
// protocols. The name of the backend ports must match the names of the
// frontend ports of the corresponding service.
//
// The Endpoints object is parsed from either an EndpointSlice (preferred) or Endpoint
// Kubernetes objects depending on the Kubernetes version.
//
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +deepequal-gen=true
// +deepequal-gen:private-method=true
type Endpoints struct {
	types.UnserializableObject
	slim_metav1.ObjectMeta

	EndpointSliceID

	// Backends is a map containing all backend IPs and ports. The key to
	// the map is the backend IP in string form. The value defines the list
	// of ports for that backend IP, plus an additional optional node name.
	// Backends map[cmtypes.AddrCluster]*Backend
	Backends map[cmtypes.AddrCluster]*Backend
}

// DeepEqual returns true if both endpoints are deep equal.
func (e *Endpoints) DeepEqual(o *Endpoints) bool {
	switch {
	case (e == nil) != (o == nil):
		return false
	case (e == nil) && (o == nil):
		return true
	}
	return e.deepEqual(o)
}

func (in *Endpoints) DeepCopyInto(out *Endpoints) {
	*out = *in
	if in.Backends != nil {
		in, out := &in.Backends, &out.Backends
		*out = make(map[cmtypes.AddrCluster]*Backend, len(*in))
		for key, val := range *in {
			var outVal *Backend
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(Backend)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
}

func (in *Endpoints) DeepCopy() *Endpoints {
	if in == nil {
		return nil
	}
	out := new(Endpoints)
	in.DeepCopyInto(out)
	return out
}

// Backend contains all ports, terminating state, and the node name of a given backend
//
// +k8s:deepcopy-gen=true
// +deepequal-gen=false
type Backend struct {
	Ports         map[loadbalancer.L4Addr][]string
	NodeName      string
	Hostname      string
	Terminating   bool
	HintsForZones []string
	Preferred     bool
	Zone          string
}

func (b *Backend) DeepEqual(other *Backend) bool {
	return maps.EqualFunc(b.Ports, other.Ports, slices.Equal) &&
		b.NodeName == other.NodeName &&
		b.Hostname == other.Hostname &&
		b.Terminating == other.Terminating &&
		slices.Equal(b.HintsForZones, other.HintsForZones) &&
		b.Preferred == other.Preferred &&
		b.Zone == other.Zone
}

func (b *Backend) ToPortConfiguration() serviceStore.PortConfiguration {
	pc := serviceStore.PortConfiguration{}
	for addr, names := range b.Ports {
		for _, name := range names {
			pc[name] = &addr
		}
	}
	return pc
}

// String returns the string representation of an endpoints resource, with
// backends and ports sorted.
func (e *Endpoints) String() string {
	if e == nil {
		return ""
	}

	backends := []string{}
	for addrCluster, be := range e.Backends {
		for port := range be.Ports {
			if be.Zone != "" {
				backends = append(backends, fmt.Sprintf("%s/%s[%s]", net.JoinHostPort(addrCluster.Addr().String(), strconv.Itoa(int(port.Port))), port.Protocol, be.Zone))
			} else {
				backends = append(backends, fmt.Sprintf("%s/%s", net.JoinHostPort(addrCluster.Addr().String(), strconv.Itoa(int(port.Port))), port.Protocol))
			}
		}
	}

	slices.Sort(backends)

	return strings.Join(backends, ",")
}

// newEndpoints returns a new Endpoints
func newEndpoints() *Endpoints {
	return &Endpoints{
		Backends: map[cmtypes.AddrCluster]*Backend{},
	}
}

// Prefixes returns the endpoint's backends as a slice of netip.Prefix.
func (e *Endpoints) Prefixes() []netip.Prefix {
	prefixes := make([]netip.Prefix, 0, len(e.Backends))
	for addrCluster := range e.Backends {
		addr := addrCluster.Addr()
		prefixes = append(prefixes, netip.PrefixFrom(addr, addr.BitLen()))
	}
	return prefixes
}

type endpointSlice interface {
	GetNamespace() string
	GetName() string
	GetLabels() map[string]string
}

// ParseEndpointSliceID parses a Kubernetes endpoints slice and returns a
// EndpointSliceID
func ParseEndpointSliceID(es endpointSlice) EndpointSliceID {
	return EndpointSliceID{
		ServiceName: loadbalancer.NewServiceName(
			es.GetNamespace(),
			es.GetLabels()[slim_discovery_v1.LabelServiceName],
		),
		EndpointSliceName: es.GetName(),
	}
}

// ParseEndpointSliceV1 parses a Kubernetes EndpointSlice resource.
// It reads ready and terminating state of endpoints in the EndpointSlice to
// return an EndpointSlice ID and a filtered list of Endpoints for service load-balancing.
func ParseEndpointSliceV1(logger *slog.Logger, ep *slim_discovery_v1.EndpointSlice) *Endpoints {
	endpoints := newEndpoints()
	endpoints.ObjectMeta = ep.ObjectMeta
	endpoints.EndpointSliceID = ParseEndpointSliceID(ep)

	// Validate AddressType before parsing. Currently, we only support IPv4 and IPv6.
	if ep.AddressType != slim_discovery_v1.AddressTypeIPv4 &&
		ep.AddressType != slim_discovery_v1.AddressTypeIPv6 {
		return endpoints
	}

	logger.Debug("Processing endpoints for EndpointSlice",
		logfields.LenEndpoints, len(ep.Endpoints),
		logfields.Name, ep.Name,
	)
	for _, sub := range ep.Endpoints {
		// ready indicates that this endpoint is prepared to receive traffic,
		// according to whatever system is managing the endpoint. A nil value
		// indicates an unknown state. In most cases consumers should interpret this
		// unknown state as ready.
		// More info: vendor/k8s.io/api/discovery/v1/types.go
		isReady := sub.Conditions.Ready == nil || *sub.Conditions.Ready
		// serving is identical to ready except that it is set regardless of the
		// terminating state of endpoints. This condition should be set to true for
		// a ready endpoint that is terminating. If nil, consumers should defer to
		// the ready condition.
		// More info: vendor/k8s.io/api/discovery/v1/types.go
		isServing := (sub.Conditions.Serving == nil && isReady) || (sub.Conditions.Serving != nil && *sub.Conditions.Serving)
		// Terminating indicates that the endpoint is getting terminated. A
		// nil values indicates an unknown state. Ready is never true when
		// an endpoint is terminating. Propagate the terminating endpoint
		// state so that we can gracefully remove those endpoints.
		// More info: vendor/k8s.io/api/discovery/v1/types.go
		isTerminating := sub.Conditions.Terminating != nil && *sub.Conditions.Terminating

		// if is not Ready allow endpoints that are Serving and Terminating
		if !isReady {

			// filter not Serving endpoints since those can not receive traffic
			if !isServing {
				logger.Debug(
					"discarding Endpoint on EndpointSlice: not Serving",
					logfields.Name, ep.Name,
				)
				continue
			}
		}

		for _, addr := range sub.Addresses {
			addrCluster, err := cmtypes.ParseAddrCluster(addr)
			if err != nil {
				logger.Info(
					"Unable to parse address for EndpointSlices",
					logfields.Error, err,
					logfields.Address, addr,
					logfields.Name, ep.Name,
				)
				continue
			}

			backend, ok := endpoints.Backends[addrCluster]
			if !ok {
				backend = &Backend{Ports: map[loadbalancer.L4Addr][]string{}}
				endpoints.Backends[addrCluster] = backend
				if sub.NodeName != nil {
					backend.NodeName = *sub.NodeName
				} else {
					if nodeName, ok := sub.DeprecatedTopology[corev1.LabelHostname]; ok {
						backend.NodeName = nodeName
					}
				}
				if sub.Hostname != nil {
					backend.Hostname = *sub.Hostname
				}
				if sub.Zone != nil {
					backend.Zone = *sub.Zone
				} else if zoneName, ok := sub.DeprecatedTopology[corev1.LabelTopologyZone]; ok {
					backend.Zone = zoneName
				}
				// If is not ready check if is serving and terminating
				if !isReady &&
					isServing && isTerminating {
					logger.Debug(
						"Endpoint address on EndpointSlice is Terminating",
						logfields.Address, addr,
						logfields.Name, ep.Name,
					)
					backend.Terminating = true
					metrics.TerminatingEndpointsEvents.Inc()
				}
			}

			for _, port := range ep.Ports {
				name, lbPort, ok := parseEndpointPortV1(port)
				if ok {
					backend.Ports[lbPort] = append(backend.Ports[lbPort], name)
				}
			}
			if sub.Hints != nil && (*sub.Hints).ForZones != nil {
				hints := (*sub.Hints).ForZones
				backend.HintsForZones = make([]string, len(hints))
				for i, hint := range hints {
					backend.HintsForZones[i] = hint.Name
				}
			}
		}
	}

	logger.Debug(
		"EndpointSlice has backends",
		logfields.LenBackends, len(endpoints.Backends),
		logfields.Name, ep.Name,
	)
	return endpoints
}

// parseEndpointPortV1 returns the port name and the port parsed as a L4Addr from
// the given port.
func parseEndpointPortV1(port slim_discovery_v1.EndpointPort) (name string, addr loadbalancer.L4Addr, ok bool) {
	proto := loadbalancer.TCP
	if port.Protocol != nil {
		switch *port.Protocol {
		case slim_corev1.ProtocolTCP:
			proto = loadbalancer.TCP
		case slim_corev1.ProtocolUDP:
			proto = loadbalancer.UDP
		case slim_corev1.ProtocolSCTP:
			proto = loadbalancer.SCTP
		default:
			return
		}
	}
	if port.Port == nil {
		return
	}
	if port.Name != nil {
		name = *port.Name
	}
	return name, loadbalancer.NewL4Addr(proto, uint16(*port.Port)), true
}
