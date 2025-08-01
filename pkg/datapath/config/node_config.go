// Code generated by dpgen. DO NOT EDIT.

// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package config

// Node is a configuration struct for a Cilium datapath object. Warning: do not
// instantiate directly! Always use [NewNode] to ensure the default values
// configured in the ELF are honored.
type Node struct {
	// Index of the interface used to connect nodes in the cluster.
	DirectRoutingDevIfindex uint32 `config:"direct_routing_dev_ifindex"`
	// Internal IPv6 router address assigned to the cilium_host interface.
	RouterIPv6 [16]byte `config:"router_ipv6"`
	// IPv4 source address used for SNAT when a Pod talks to itself over a Service.
	ServiceLoopbackIPv4 [4]byte `config:"service_loopback_ipv4"`
	// Length of payload to capture when tracing native packets.
	TracePayloadLen uint32 `config:"trace_payload_len"`
	// Length of payload to capture when tracing overlay packets.
	TracePayloadLenOverlay uint32 `config:"trace_payload_len_overlay"`
}

func NewNode() *Node {
	return &Node{0x0,
		[16]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		[4]byte{0x0, 0x0, 0x0, 0x0}, 0x0, 0x0}
}
