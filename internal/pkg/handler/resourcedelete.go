type ResourceDeleteHandler struct {
	Resource   interface{}
	Collectors metrics.Collectors
	Recorder   record.EventRecorder
}

// Handle processes the resource deletion handler
func (r ResourceDeleteHandler) Handle() error {
	deleted := r.Resource
	resourceType := ""

	objectMeta, err := getObjectMetaForResource(deleted)

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
	if objectMeta.Namespace != "" {
		ns, err := clients.KubernetesClient.CoreV1().Namespaces().Get(context.TODO(), objectMeta.Namespace, metav1.GetOptions{})
		if err == nil && ns != nil && ns.Annotations != nil {
			if pauseValue, found := ns.Annotations[options.PauseReloaderAnnotation]; found {
				isPaused, err := strconv.ParseBool(pauseValue)
				if err == nil && isPaused {
					logrus.Infof("Namespace %s is paused, skipping reload for resource %s", objectMeta.Namespace, objectMeta.Name)
					return nil
				}
			}
		}
	}

	switch resource := deleted.(type) {
	case *v1.ConfigMap:
		if !needsReload(resource, objectMeta, r.Recorder) {
			logrus.Infof("Delete update ignored for ConfigMap: %s", resource.Name)
			return nil
		}
		resourceType = "configmap"
	case *v1.Secret:
		if !needsReload(resource, objectMeta, r.Recorder) {
			logrus.Infof("Delete update ignored for Secret: %s", resource.Name)
			return nil
		}
		resourceType = "secret"
	}

	resourceInfo, err := r.getResourceInfo(deleted)

	if err != nil {
		return err
	}

	if err := r.resourceDeleted(resourceInfo, objectMeta, resourceType); err != nil {
		return err
	}

	return nil
}
