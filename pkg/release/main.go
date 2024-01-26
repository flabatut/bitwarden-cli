package release

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger.io/dagger"
)

type Workflow struct {
	Client           *dagger.Client
	ProjectNamespace string
	RegistryPassword *dagger.Secret
}

func (w *Workflow) Release(ctx context.Context, releaseVersion string) (*dagger.Container, error) {
	fmt.Println("Releasing with Dagger")

	var ghRepo = fmt.Sprintf("github.com/%s", w.ProjectNamespace)

	releaser := w.Client.Container().
		From("alpine:latest").
		// WithEnvVariable("GH_DEBUG", "api").
		WithEnvVariable("GH_REPO", ghRepo).
		WithSecretVariable("GH_TOKEN", w.RegistryPassword).
		WithExec([]string{"apk", "add", "github-cli"})

	// https://docs.dagger.io/cookbook/#invalidate-cache
	_, err := releaser.WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{
			"gh", "release", "view", releaseVersion,
		}).Stdout(ctx)
	if err != nil {
		if strings.HasSuffix(err.Error(), "release not found") {
			_, err = releaser.WithEnvVariable("CACHEBUSTER", time.Now().String()).
				WithExec([]string{
					"gh", "release", "create", releaseVersion, "-t", releaseVersion, "--generate-notes",
				}).Stdout(ctx)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return releaser, nil
}
