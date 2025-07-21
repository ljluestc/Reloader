package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stakater/Reloader/internal/pkg/annotations"
	"github.com/stakater/Reloader/internal/pkg/metrics"
	"github.com/stakater/Reloader/internal/pkg/options"
	"github.com/stakater/Reloader/pkg/kube"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"strconv"
)

// ResourceUpdatedHandler contains update update which need to be processed
type ResourceUpdatedHandler struct {
	Resource    interface{}
	OldResource interface{}
	Collectors  metrics.Collectors
	Recorder    record.EventRecorder
}

// Handle processes the update handler
func (r ResourceUpdatedHandler) Handle() error {
	updated := r.Resource

	oldResource := r.OldResource
	newResource := r.Resource
	resourceType := ""

	objectMeta, oldObjectMeta, err := getObjectMetaForHandleResource(oldResource, newResource)

	if err != nil {
		return err
	}

	clients := kube.GetClients()

	// Check if the resource or namespace is paused
	if pauseValue, found := objectMeta.Annotations[annotations.PauseReloaderAnnotation]; found {
		isPaused, err := strconv.ParseBool(pauseValue)
		if err == nil && isPaused {
			logrus.Infof("Resource %s/%s is paused, skipping reload", objectMeta.Namespace, objectMeta.Name)
			return nil
		}
	}

	// Check if the namespace is paused
	if objectMeta.Namespace != "" {
		ns, err := clients.KubernetesClient.CoreV1().Namespaces().Get(context.TODO(), objectMeta.Namespace, metav1.GetOptions{})
		if err == nil && ns != nil && ns.Annotations != nil {
			if pauseValue, found := ns.Annotations[annotations.PauseReloaderAnnotation]; found {
				isPaused, err := strconv.ParseBool(pauseValue)
				if err == nil && isPaused {
					logrus.Infof("Namespace %s is paused, skipping reload for resource %s", objectMeta.Namespace, objectMeta.Name)
					return nil
				}
			}
		}
	}

	switch resource := updated.(type) {
	case *v1.ConfigMap:
		config, oldConfig, err := getConfigmapConfig(oldResource, newResource)
		if err != nil {
			return err
		}
		createInCluster, err := checkResourceUpdate(config, oldConfig)
		if err != nil {
			return err
		}
		if !createInCluster && !needsReload(resource, objectMeta, r.Recorder) {
			logrus.Infof("Add update ignored for ConfigMap: %s", resource.Name)
			return nil
		}
		resourceType = "configmap"
	case *v1.Secret:
		config, oldConfig, err := getSecretConfig(oldResource, newResource)
		if err != nil {
			return err
		}
		createInCluster, err := checkResourceUpdate(config, oldConfig)
		if err != nil {
			return err
		}
		if !createInCluster && !needsReload(resource, objectMeta, r.Recorder) {
			logrus.Infof("Add update ignored for Secret: %s", resource.Name)
			return nil
		}
		resourceType = "secret"
	}

	resourceInfo, err := r.getResourceInfo(updated)

	if err != nil {
		return err
	}

	if err := r.resourceUpdated(resourceInfo, objectMeta, oldObjectMeta, resourceType); err != nil {
		return err
	}

	return nil
}
