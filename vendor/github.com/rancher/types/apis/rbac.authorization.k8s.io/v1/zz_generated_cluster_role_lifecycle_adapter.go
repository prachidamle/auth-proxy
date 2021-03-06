package v1

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ClusterRoleLifecycle interface {
	Create(obj *v1.ClusterRole) error
	Remove(obj *v1.ClusterRole) error
	Updated(obj *v1.ClusterRole) error
}

type clusterRoleLifecycleAdapter struct {
	lifecycle ClusterRoleLifecycle
}

func (w *clusterRoleLifecycleAdapter) Create(obj runtime.Object) error {
	return w.lifecycle.Create(obj.(*v1.ClusterRole))
}

func (w *clusterRoleLifecycleAdapter) Finalize(obj runtime.Object) error {
	return w.lifecycle.Remove(obj.(*v1.ClusterRole))
}

func (w *clusterRoleLifecycleAdapter) Updated(obj runtime.Object) error {
	return w.lifecycle.Updated(obj.(*v1.ClusterRole))
}

func NewClusterRoleLifecycleAdapter(name string, client ClusterRoleInterface, l ClusterRoleLifecycle) ClusterRoleHandlerFunc {
	adapter := &clusterRoleLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, adapter, client.ObjectClient())
	return func(key string, obj *v1.ClusterRole) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
