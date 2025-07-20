package kube

import (
	"context"
	"fmt"
	"os"

	argorollout "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	appsclient "github.com/openshift/client-go/apps/clientset/versioned"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Clients struct exposes interfaces for kubernetes as well as openshift if available
type Clients struct {
	KubernetesClient    kubernetes.Interface
	KubeConfig          *rest.Config
	OpenshiftAppsClient appsclient.Interface
	ArgoRolloutClient   argorollout.Interface
}

var (
	// IsOpenshift is true if environment is Openshift, it is false if environment is Kubernetes
	IsOpenshift bool
)

// GetClient is a function variable that can be overridden in tests
var GetClient = func() (interface{}, error) {
	clients, err := GetClientInCluster()
	if err != nil {
		// If in-cluster config fails, try out-of-cluster
		clients, err = GetClientOutOfCluster()
		if err != nil {
			return nil, fmt.Errorf("unable to create Kubernetes client error = %v", err)
		}
	}
	return clients, nil
}

func init() {
	IsOpenshift = isOpenshift()
}

// GetClientOutOfCluster returns a k8s clientset and kubeconfig to the cluster from outside
func GetClientOutOfCluster() (*Clients, error) {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Clients{
		KubernetesClient: client,
		KubeConfig:       config,
	}, nil
}

// GetClientInCluster returns a k8s clientset to the cluster from within
func GetClientInCluster() (*Clients, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Clients{
		KubernetesClient: client,
		KubeConfig:       config,
	}, nil
}

// GetClients returns a `Clients` object containing both openshift and kubernetes clients with an openshift identifier
func GetClients() Clients {
	clients, err := GetClient()
	if err != nil {
		logrus.Fatalf("Unable to create Kubernetes client error = %v", err)
	}

	k8sClient := clients.(*Clients).KubernetesClient
	if k8sClient == nil {
		logrus.Fatalf("Unable to get Kubernetes client")
	}

	var appsClient appsclient.Interface

	if IsOpenshift {
		_, err := getConfig()
		if err != nil {
			logrus.Warnf("Unable to create Openshift Apps client error = %v", err)
		} else {
			appsClient, err = appsclient.NewForConfig(clients.(*Clients).KubeConfig)
			if err != nil {
				logrus.Warnf("Unable to create Openshift Apps client error = %v", err)
			}
		}
	}

	var rolloutClient *argorollout.Clientset

	rolloutClient, err = GetArgoRolloutClient()
	if err != nil {
		logrus.Warnf("Unable to create ArgoRollout client error = %v", err)
	}

	return Clients{
		KubernetesClient:    k8sClient,
		OpenshiftAppsClient: appsClient,
		ArgoRolloutClient:   rolloutClient,
	}
}

func GetArgoRolloutClient() (*argorollout.Clientset, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	return argorollout.NewForConfig(config)
}

func isOpenshift() bool {
	cfg, err := getConfig()
	if err != nil {
		logrus.Fatalf("Unable to get Kubernetes config error = %v", err)
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logrus.Fatalf("Unable to create Kubernetes client error = %v", err)
	}
	_, err = client.RESTClient().Get().AbsPath("/apis/project.openshift.io").Do(context.TODO()).Raw()
	if err == nil {
		logrus.Info("Environment: Openshift")
		return true
	}
	logrus.Info("Environment: Kubernetes")
	return false
}

// GetOpenshiftAppsClient returns an Openshift Client that can query on Apps
func GetOpenshiftAppsClient() (*appsclient.Clientset, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	return appsclient.NewForConfig(config)
}

// GetKubernetesClient gets the client for k8s, if ~/.kube/config exists so get that config else incluster config
func GetKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func getConfig() (*rest.Config, error) {
	var config *rest.Config
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	//If file exists so use that config settings
	if _, err := os.Stat(kubeconfigPath); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, err
		}
	} else { //Use Incluster Configuration
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
