// +build !ignore_autogenerated

/*
Copyright 2021 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BrokerCell) DeepCopyInto(out *BrokerCell) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BrokerCell.
func (in *BrokerCell) DeepCopy() *BrokerCell {
	if in == nil {
		return nil
	}
	out := new(BrokerCell)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BrokerCell) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BrokerCellList) DeepCopyInto(out *BrokerCellList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BrokerCell, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BrokerCellList.
func (in *BrokerCellList) DeepCopy() *BrokerCellList {
	if in == nil {
		return nil
	}
	out := new(BrokerCellList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BrokerCellList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BrokerCellSpec) DeepCopyInto(out *BrokerCellSpec) {
	*out = *in
	in.Components.DeepCopyInto(&out.Components)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BrokerCellSpec.
func (in *BrokerCellSpec) DeepCopy() *BrokerCellSpec {
	if in == nil {
		return nil
	}
	out := new(BrokerCellSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BrokerCellStatus) DeepCopyInto(out *BrokerCellStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BrokerCellStatus.
func (in *BrokerCellStatus) DeepCopy() *BrokerCellStatus {
	if in == nil {
		return nil
	}
	out := new(BrokerCellStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComponentParameters) DeepCopyInto(out *ComponentParameters) {
	*out = *in
	if in.AvgCPUUtilization != nil {
		in, out := &in.AvgCPUUtilization, &out.AvgCPUUtilization
		*out = new(int32)
		**out = **in
	}
	if in.AvgMemoryUsage != nil {
		in, out := &in.AvgMemoryUsage, &out.AvgMemoryUsage
		*out = new(string)
		**out = **in
	}
	if in.MinReplicas != nil {
		in, out := &in.MinReplicas, &out.MinReplicas
		*out = new(int32)
		**out = **in
	}
	if in.MaxReplicas != nil {
		in, out := &in.MaxReplicas, &out.MaxReplicas
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComponentParameters.
func (in *ComponentParameters) DeepCopy() *ComponentParameters {
	if in == nil {
		return nil
	}
	out := new(ComponentParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComponentsParametersSpec) DeepCopyInto(out *ComponentsParametersSpec) {
	*out = *in
	if in.Fanout != nil {
		in, out := &in.Fanout, &out.Fanout
		*out = new(ComponentParameters)
		(*in).DeepCopyInto(*out)
	}
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(ComponentParameters)
		(*in).DeepCopyInto(*out)
	}
	if in.Retry != nil {
		in, out := &in.Retry, &out.Retry
		*out = new(ComponentParameters)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComponentsParametersSpec.
func (in *ComponentsParametersSpec) DeepCopy() *ComponentsParametersSpec {
	if in == nil {
		return nil
	}
	out := new(ComponentsParametersSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSpecification) DeepCopyInto(out *ResourceSpecification) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSpecification.
func (in *ResourceSpecification) DeepCopy() *ResourceSpecification {
	if in == nil {
		return nil
	}
	out := new(ResourceSpecification)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SystemResource) DeepCopyInto(out *SystemResource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SystemResource.
func (in *SystemResource) DeepCopy() *SystemResource {
	if in == nil {
		return nil
	}
	out := new(SystemResource)
	in.DeepCopyInto(out)
	return out
}