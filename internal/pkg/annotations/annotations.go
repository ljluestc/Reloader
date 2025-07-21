package annotations

// Annotation keys for auto reloaders
const (
	ReloaderAutoAnnotation             = "reloader.stakater.com/auto"
	ConfigmapReloaderAutoAnnotation    = "configmap.reloader.stakater.com/auto"
	SecretReloaderAutoAnnotation       = "secret.reloader.stakater.com/auto"
	AutoSearchAnnotation               = "reloader.stakater.com/search"
	SearchMatchAnnotation              = "reloader.stakater.com/match"
	ConfigmapReloaderAnnotation        = "configmap.reloader.stakater.com/reload"
	SecretReloaderAnnotation           = "secret.reloader.stakater.com/reload"
	ReloaderAnnotation                 = "reloader.stakater.com/reload"
	ReloadStrategyAnnotation           = "reloader.stakater.com/reload-strategy"
	ReloadOnRestartStrategyValue       = "restart"
	ReloadOnRolloutStrategyValue       = "rollout"
	ReloaderIgnoreAnnotation           = "reloader.stakater.com/ignore"
	ConfigmapExcludeReloaderAnnotation = "configmap.reloader.stakater.com/exclude"
	SecretExcludeReloaderAnnotation    = "secret.reloader.stakater.com/exclude"
	// PauseReloaderAnnotation is used to pause reloading for a resource or namespace
	PauseReloaderAnnotation            = "reloader.stakater.com/pause"
)
