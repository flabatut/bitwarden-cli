package publish

import (
	"dagger.io/dagger"
)

type PublishWorkflow struct {
	Client           *dagger.Client
	PlatformVariants []*dagger.Container
	PublishAddr      string
	// ReleaseVersion       string
	// BuilderNodeJSVersion string
	// RunnerEntryPointPath string
	// BuilderWorkDir       string
	// BuilderImage         string
	// RunnerImage          string
	// BuilderPlatforms     []dagger.Platform
}
