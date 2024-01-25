package image

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

type Workflow struct {
	Client           *dagger.Client
	ReleaseVersion   string
	RegistryFQDN     string
	ProjectNamespace string
	RegistryUsername string
	RegistryPassword *dagger.Secret
}

func (w *Workflow) Publish(ctx context.Context, platformVariants []*dagger.Container) error {
	fmt.Println("Publishing with Dagger")
	// TODO: validate registry URL
	// TODO: support image name option
	publishAddress := fmt.Sprintf("%s/%s:%s", w.RegistryFQDN, w.ProjectNamespace, w.ReleaseVersion)
	imageDigest, err := w.Client.Container().
		WithRegistryAuth(w.RegistryFQDN, w.RegistryUsername, w.RegistryPassword).
		Publish(ctx, publishAddress, dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
		})
	if err != nil {
		return err
	}
	fmt.Println("published multi-platform image with digest", imageDigest)

	return nil
}
