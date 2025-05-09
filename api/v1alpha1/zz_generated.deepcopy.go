//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupDaemon) DeepCopyInto(out *BackupDaemon) {
	*out = *in
	in.Affinity.DeepCopyInto(&out.Affinity)
	in.BackupStorage.DeepCopyInto(&out.BackupStorage)
	in.Resources.DeepCopyInto(&out.Resources)
	in.SecurityContext.DeepCopyInto(&out.SecurityContext)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupDaemon.
func (in *BackupDaemon) DeepCopy() *BackupDaemon {
	if in == nil {
		return nil
	}
	out := new(BackupDaemon)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupDaemonStatus) DeepCopyInto(out *BackupDaemonStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupDaemonStatus.
func (in *BackupDaemonStatus) DeepCopy() *BackupDaemonStatus {
	if in == nil {
		return nil
	}
	out := new(BackupDaemonStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupStorage) DeepCopyInto(out *BackupStorage) {
	*out = *in
	if in.StorageClass != nil {
		in, out := &in.StorageClass, &out.StorageClass
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupStorage.
func (in *BackupStorage) DeepCopy() *BackupStorage {
	if in == nil {
		return nil
	}
	out := new(BackupStorage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Monitoring) DeepCopyInto(out *Monitoring) {
	*out = *in
	in.Affinity.DeepCopyInto(&out.Affinity)
	in.Resources.DeepCopyInto(&out.Resources)
	in.SecurityContext.DeepCopyInto(&out.SecurityContext)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Monitoring.
func (in *Monitoring) DeepCopy() *Monitoring {
	if in == nil {
		return nil
	}
	out := new(Monitoring)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitoringStatus) DeepCopyInto(out *MonitoringStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringStatus.
func (in *MonitoringStatus) DeepCopy() *MonitoringStatus {
	if in == nil {
		return nil
	}
	out := new(MonitoringStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SnapshotStorage) DeepCopyInto(out *SnapshotStorage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SnapshotStorage.
func (in *SnapshotStorage) DeepCopy() *SnapshotStorage {
	if in == nil {
		return nil
	}
	out := new(SnapshotStorage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Storage) DeepCopyInto(out *Storage) {
	*out = *in
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ClassName != nil {
		in, out := &in.ClassName, &out.ClassName
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Storage.
func (in *Storage) DeepCopy() *Storage {
	if in == nil {
		return nil
	}
	out := new(Storage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZooKeeper) DeepCopyInto(out *ZooKeeper) {
	*out = *in
	in.Affinity.DeepCopyInto(&out.Affinity)
	in.Storage.DeepCopyInto(&out.Storage)
	out.SnapshotStorage = in.SnapshotStorage
	in.Resources.DeepCopyInto(&out.Resources)
	in.SecurityContext.DeepCopyInto(&out.SecurityContext)
	if in.EnvironmentVariables != nil {
		in, out := &in.EnvironmentVariables, &out.EnvironmentVariables
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZooKeeper.
func (in *ZooKeeper) DeepCopy() *ZooKeeper {
	if in == nil {
		return nil
	}
	out := new(ZooKeeper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZooKeeperService) DeepCopyInto(out *ZooKeeperService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZooKeeperService.
func (in *ZooKeeperService) DeepCopy() *ZooKeeperService {
	if in == nil {
		return nil
	}
	out := new(ZooKeeperService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ZooKeeperService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZooKeeperServiceList) DeepCopyInto(out *ZooKeeperServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ZooKeeperService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZooKeeperServiceList.
func (in *ZooKeeperServiceList) DeepCopy() *ZooKeeperServiceList {
	if in == nil {
		return nil
	}
	out := new(ZooKeeperServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ZooKeeperServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZooKeeperServiceSpec) DeepCopyInto(out *ZooKeeperServiceSpec) {
	*out = *in
	if in.ZooKeeper != nil {
		in, out := &in.ZooKeeper, &out.ZooKeeper
		*out = new(ZooKeeper)
		(*in).DeepCopyInto(*out)
	}
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(Monitoring)
		(*in).DeepCopyInto(*out)
	}
	if in.BackupDaemon != nil {
		in, out := &in.BackupDaemon, &out.BackupDaemon
		*out = new(BackupDaemon)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZooKeeperServiceSpec.
func (in *ZooKeeperServiceSpec) DeepCopy() *ZooKeeperServiceSpec {
	if in == nil {
		return nil
	}
	out := new(ZooKeeperServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZooKeeperServiceStatus) DeepCopyInto(out *ZooKeeperServiceStatus) {
	*out = *in
	in.ZooKeeperStatus.DeepCopyInto(&out.ZooKeeperStatus)
	in.MonitoringStatus.DeepCopyInto(&out.MonitoringStatus)
	in.BackupDaemonStatus.DeepCopyInto(&out.BackupDaemonStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZooKeeperServiceStatus.
func (in *ZooKeeperServiceStatus) DeepCopy() *ZooKeeperServiceStatus {
	if in == nil {
		return nil
	}
	out := new(ZooKeeperServiceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZooKeeperStatus) DeepCopyInto(out *ZooKeeperStatus) {
	*out = *in
	if in.Servers != nil {
		in, out := &in.Servers, &out.Servers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZooKeeperStatus.
func (in *ZooKeeperStatus) DeepCopy() *ZooKeeperStatus {
	if in == nil {
		return nil
	}
	out := new(ZooKeeperStatus)
	in.DeepCopyInto(out)
	return out
}
