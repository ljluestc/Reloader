package kube

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceMapTypes are resources from where changes are going to be detected
var ResourceMapTypes = map[string]runtime.Object{
	"configMaps": &v1.ConfigMap{},
	"secrets":    &v1.Secret{},
	"namespaces": &v1.Namespace{},
}
