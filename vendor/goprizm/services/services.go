// services pkg provides helpers to discover/access prizm services and resources.
// Configs will be available to apps as environment variables. Environment variables
// will populated as Kubernetes ConfigMap(https://kubernetes.io/docs/tasks/configure-pod-container/configmap/).
// For dev testing purpose environ vars can also be manauly set.
//
// ConfigMap as Volume - https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#populate-a-volume-with-data-stored-in-a-configmap).
// ConfigMap can also be mounted as volume at /var/prizm. This can be considered in future.
//
// If other service discovery mechanisms are used in future helper funcs can modified
// appropriately so that application services need not change.
package services

import (
	"path/filepath"
)

var (
	Var = "/var/aruba/prizm"
	Etc = filepath.Join(Var, "/etc")
)
