// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package cidr

import (
	"bytes"
	"fmt"
	"net"
	"slices"
)

// NewCIDR returns a new CIDR using a net.IPNet
func NewCIDR(ipnet *net.IPNet) *CIDR {
	if ipnet == nil {
		return nil
	}

	return &CIDR{ipnet}
}

func NewCIDRSlice(ipnets []*net.IPNet) []*CIDR {
	if ipnets == nil {
		return nil
	}

	cidrs := make([]*CIDR, len(ipnets))
	for i, ipnet := range ipnets {
		cidrs[i] = NewCIDR(ipnet)
	}
	return cidrs
}

func CIDRsToIPNets(cidrs []*CIDR) []*net.IPNet {
	if cidrs == nil {
		return nil
	}

	ipnets := make([]*net.IPNet, len(cidrs))
	for i, cidr := range cidrs {
		ipnets[i] = cidr.IPNet
	}
	return ipnets
}

// CIDR is a network CIDR representation based on net.IPNet
type CIDR struct {
	*net.IPNet
}

func (c *CIDR) String() string {
	if c == nil {
		var n *net.IPNet
		return n.String()
	}
	return c.IPNet.String()
}

// DeepEqual is an deepequal function, deeply comparing the receiver with other.
// in must be non-nil.
func (in *CIDR) DeepEqual(other *CIDR) bool {
	if other == nil {
		return false
	}

	if (in.IPNet == nil) != (other.IPNet == nil) {
		return false
	} else if in.IPNet != nil {
		if !in.IPNet.IP.Equal(other.IPNet.IP) {
			return false
		}
		inOnes, inBits := in.IPNet.Mask.Size()
		otherOnes, otherBits := other.IPNet.Mask.Size()
		return inOnes == otherOnes && inBits == otherBits
	}

	return true
}

// DeepCopy creates a deep copy of a CIDR
func (n *CIDR) DeepCopy() *CIDR {
	if n == nil {
		return nil
	}
	out := new(CIDR)
	n.DeepCopyInto(out)
	return out
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CIDR) DeepCopyInto(out *CIDR) {
	*out = *in
	if in.IPNet == nil {
		return
	}
	out.IPNet = new(net.IPNet)
	*out.IPNet = *in.IPNet
	if in.IPNet.IP != nil {
		in, out := &in.IPNet.IP, &out.IPNet.IP
		*out = make(net.IP, len(*in))
		copy(*out, *in)
	}
	if in.IPNet.Mask != nil {
		in, out := &in.IPNet.Mask, &out.IPNet.Mask
		*out = make(net.IPMask, len(*in))
		copy(*out, *in)
	}
}

// Equal returns true if the receiver's CIDR equals the other CIDR.
func (n *CIDR) Equal(o *CIDR) bool {
	if n == nil || o == nil {
		return n == o
	}
	return Equal(n.IPNet, o.IPNet)
}

// Equal returns true if the n and o net.IPNet CIDRs are Equal.
func Equal(n, o *net.IPNet) bool {
	if n == nil || o == nil {
		return n == o
	}
	if n == o {
		return true
	}
	return n.IP.Equal(o.IP) &&
		bytes.Equal(n.Mask, o.Mask)
}

// ZeroNet generates a zero net.IPNet object for the given address family
func ZeroNet(family int) *net.IPNet {
	switch family {
	case FAMILY_V4:
		return &net.IPNet{
			IP:   net.IPv4zero,
			Mask: net.CIDRMask(0, 8*net.IPv4len),
		}
	case FAMILY_V6:
		return &net.IPNet{
			IP:   net.IPv6zero,
			Mask: net.CIDRMask(0, 8*net.IPv6len),
		}
	}
	return nil
}

// ContainsAll returns true if 'ipNets1' contains all net.IPNet of 'ipNets2'
func ContainsAll(ipNets1, ipNets2 []*net.IPNet) bool {
	for _, n2 := range ipNets2 {
		if !slices.ContainsFunc(ipNets1, func(n1 *net.IPNet) bool {
			return Equal(n2, n1)
		}) {
			return false
		}
	}
	return true
}

// ParseCIDR parses the CIDR string using net.ParseCIDR
func ParseCIDR(str string) (*CIDR, error) {
	_, ipnet, err := net.ParseCIDR(str)
	if err != nil {
		return nil, err
	}
	return NewCIDR(ipnet), nil
}

// MustParseCIDR parses the CIDR string using net.ParseCIDR and panics if the
// CIDR cannot be parsed
func MustParseCIDR(str string) *CIDR {
	c, err := ParseCIDR(str)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse CIDR '%s': %s", str, err))
	}
	return c
}
