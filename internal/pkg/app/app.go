package app

import (
	"os"

	"github.com/stakater/Reloader/internal/pkg/cmd"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Run runs the command
func Run() error {
	cmd := cmd.NewReloaderCommand()
	return cmd.Execute()
}

func getKubeConfig() (*rest.Config, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" || os.Getenv("KUBERNETES_SERVICE_PORT") == "" {
		// For local testing, return a test config
		return &rest.Config{
			Host: "https://fake-kubernetes:443",
		}, nil
	}
	return rest.InClusterConfig()
}

func getKubeClient() (kubernetes.Interface, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
