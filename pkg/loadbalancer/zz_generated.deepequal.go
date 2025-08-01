//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

// Code generated by deepequal-gen. DO NOT EDIT.

package loadbalancer

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *Config) DeepEqual(other *Config) bool {
	if other == nil {
		return false
	}

	if !in.UserConfig.DeepEqual(&other.UserConfig) {
		return false
	}

	if in.NodePortMin != other.NodePortMin {
		return false
	}
	if in.NodePortMax != other.NodePortMax {
		return false
	}

	return true
}

// deepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *L4Addr) deepEqual(other *L4Addr) bool {
	if other == nil {
		return false
	}

	if in.Protocol != other.Protocol {
		return false
	}
	if in.Port != other.Port {
		return false
	}

	return true
}

// DeepEqual is an autogenerated deepequal function, deeply comparing the
// receiver with other. in must be non-nil.
func (in *UserConfig) DeepEqual(other *UserConfig) bool {
	if other == nil {
		return false
	}

	if in.RetryBackoffMin != other.RetryBackoffMin {
		return false
	}
	if in.RetryBackoffMax != other.RetryBackoffMax {
		return false
	}
	if in.LBMapEntries != other.LBMapEntries {
		return false
	}
	if in.LBServiceMapEntries != other.LBServiceMapEntries {
		return false
	}
	if in.LBBackendMapEntries != other.LBBackendMapEntries {
		return false
	}
	if in.LBRevNatEntries != other.LBRevNatEntries {
		return false
	}
	if in.LBAffinityMapEntries != other.LBAffinityMapEntries {
		return false
	}
	if in.LBSourceRangeAllTypes != other.LBSourceRangeAllTypes {
		return false
	}
	if in.LBSourceRangeMapEntries != other.LBSourceRangeMapEntries {
		return false
	}
	if in.LBMaglevMapEntries != other.LBMaglevMapEntries {
		return false
	}
	if in.LBSockRevNatEntries != other.LBSockRevNatEntries {
		return false
	}
	if ((in.NodePortRange != nil) && (other.NodePortRange != nil)) || ((in.NodePortRange == nil) != (other.NodePortRange == nil)) {
		in, other := &in.NodePortRange, &other.NodePortRange
		if other == nil {
			return false
		}

		if len(*in) != len(*other) {
			return false
		} else {
			for i, inElement := range *in {
				if inElement != (*other)[i] {
					return false
				}
			}
		}
	}

	if in.LBMode != other.LBMode {
		return false
	}
	if in.LBModeAnnotation != other.LBModeAnnotation {
		return false
	}
	if in.LBAlgorithm != other.LBAlgorithm {
		return false
	}
	if in.DSRDispatch != other.DSRDispatch {
		return false
	}
	if in.ExternalClusterIP != other.ExternalClusterIP {
		return false
	}
	if in.AlgorithmAnnotation != other.AlgorithmAnnotation {
		return false
	}
	if in.EnableHealthCheckNodePort != other.EnableHealthCheckNodePort {
		return false
	}
	if in.LBPressureMetricsInterval != other.LBPressureMetricsInterval {
		return false
	}
	if in.LBSockTerminateAllProtos != other.LBSockTerminateAllProtos {
		return false
	}
	if in.EnableServiceTopology != other.EnableServiceTopology {
		return false
	}

	return true
}
