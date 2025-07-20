package handler

import (
	"github.com/stakater/Reloader/internal/pkg/util"
	"k8s.io/client-go/kubernetes"
)

// ResourceHandler handles the creation and update of resources
type ResourceHandler interface {
	Handle() error
	GetConfig() (util.Config, string)
}

type Handler struct {
	client kubernetes.Interface
}

// NewHandler creates a new handler with the given Kubernetes client
func NewHandler(client kubernetes.Interface) *Handler {
	return &Handler{
		client: client,
	}
}
