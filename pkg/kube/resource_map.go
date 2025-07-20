package kube

// ResourceMap is a map of resource types supported by Reloader
var ResourceMap = map[string]string{
	"configMaps": "true",
	"secrets":    "true",
}
