// Code generated by dpgen. DO NOT EDIT.

// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package config

// BPFHost is a configuration struct for a Cilium datapath object. Warning: do
// not instantiate directly! Always use [NewBPFHost] to ensure the default
// values configured in the ELF are honored.
type BPFHost struct {
	// MTU of the device the bpf program is attached to (default: MTU set in
	// node_config.h by agent).
	DeviceMTU uint16 `config:"device_mtu"`
	// Pass traffic with extended IP protocols.
	EnableExtendedIPProtocols bool `config:"enable_extended_ip_protocols"`
	// Length of the Ethernet header on this device. May be set to zero on L2-less
	// devices. (default __ETH_HLEN).
	EthHeaderLength uint8 `config:"eth_header_length"`
	// The host endpoint's security ID.
	HostEpID uint16 `config:"host_ep_id"`
	// Ifindex of the interface the bpf program is attached to.
	InterfaceIfindex uint32 `config:"interface_ifindex"`
	// MAC address of the interface the bpf program is attached to.
	InterfaceMAC [8]byte `config:"interface_mac"`
	// Masquerade address for IPv4 traffic.
	NATIPv4Masquerade [4]byte `config:"nat_ipv4_masquerade"`
	// Masquerade address for IPv6 traffic.
	NATIPv6Masquerade [16]byte `config:"nat_ipv6_masquerade"`
	// Pull security context from IP cache.
	SecctxFromIPCache bool `config:"secctx_from_ipcache"`
	// The endpoint's security label.
	SecurityLabel uint32 `config:"security_label"`

	Node
}

func NewBPFHost(node Node) *BPFHost {
	return &BPFHost{0x5dc, false, 0xe, 0x0, 0x0, [8]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		[4]byte{0x0, 0x0, 0x0, 0x0},
		[16]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		false, 0x0, node}
}
