package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stakater/Reloader/internal/pkg/metrics"
	"github.com/stakater/Reloader/internal/pkg/options"

	"github.com/stakater/Reloader/pkg/kube"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"strconv"
)

// ResourceCreatedHandler contains resource create which need to be processed
type ResourceCreatedHandler struct {
	Resource   interface{}
	Collectors metrics.Collectors
	Recorder   record.EventRecorder
}

// Handle processes the resource creation handler
func (r ResourceCreatedHandler) Handle() error {
	created := r.Resource
	resourceType := ""

	objectMeta, err := getObjectMetaForResource(created)

	if err != nil {
		return err
	}

	clients := kube.GetClients()

	// Check if the resource is paused
	if pauseValue, found := objectMeta.Annotations[options.PauseReloaderAnnotation]; found {
		isPaused, err := strconv.ParseBool(pauseValue)
		if err == nil && isPaused {
			logrus.Infof("Resource %s/%s is paused, skipping reload", objectMeta.Namespace, objectMeta.Name)
			return nil
		}
	}

	// Check if the namespace is paused
	if objectMeta.GetNamespace() != "" {
		ns, err := clients.KubernetesClient.CoreV1().Namespaces().Get(context.TODO(), objectMeta.GetNamespace(), metav1.GetOptions{})
		if err == nil && ns != nil && ns.Annotations != nil {
			if pauseValue, found := ns.Annotations[options.PauseReloaderAnnotation]; found {
				isPaused, err := strconv.ParseBool(pauseValue)
				if err == nil && isPaused {
					logrus.Infof("Namespace %s is paused, skipping reload for resource %s", objectMeta.GetNamespace(), objectMeta.Name)
					return nil
				}
			}
		}
	}

	switch resource := created.(type) {
	case *v1.ConfigMap:
		if !needsReload(resource, objectMeta, r.Recorder) {
			logrus.Infof("Add update ignored for ConfigMap: %s", resource.Name)
			return nil
		}
		resourceType = "configmap"
	case *v1.Secret:
		if !needsReload(resource, objectMeta, r.Recorder) {
			logrus.Infof("Add update ignored for Secret: %s", resource.Name)
			return nil
		}
		resourceType = "secret"
	}

	resourceInfo, err := r.getResourceInfo(created)

	if err != nil {
		return err
	}

	if err := r.resourceCreated(resourceInfo, objectMeta, resourceType); err != nil {
		return err
	}

	return nil
}
