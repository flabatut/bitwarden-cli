package artifact

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// TODO factorize
const containerDistDir = "/dist"

type Workflow struct {
	Client           *dagger.Client
	ReleaseVersion   string
	RegistryFQDN     string
	ProjectNamespace string
	// RegistryUsername string
	RegistryPassword  *dagger.Secret
	PlatformsVariants []dagger.Platform
}

func (w *Workflow) Publish(ctx context.Context, artifacts *dagger.Directory) error {
	fmt.Println("Publishing with Dagger")

	var ghRepo = fmt.Sprintf("github.com/%s", w.ProjectNamespace)

	publisher := w.Client.Container().
		From("alpine:latest").
		WithMountedDirectory(containerDistDir, artifacts).
		WithEnvVariable("GH_DEBUG", "api").
		WithEnvVariable("GH_REPO", ghRepo).
		WithSecretVariable("GH_TOKEN", w.RegistryPassword).
		WithWorkdir(containerDistDir).
		WithExec([]string{"apk", "add", "github-cli"})

	_, err := publisher.WithExec([]string{
		"gh", "release", "view", w.ReleaseVersion,
	}).Stdout(ctx)
	if err != nil {
		if strings.HasSuffix(err.Error(), "release not found") {
			_, err = publisher.WithExec([]string{
				"gh", "release", "create", w.ReleaseVersion, "-t", w.ReleaseVersion, "--generate-notes",
			}).Stdout(ctx)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	for _, platform := range w.PlatformsVariants {
		// extract os/arch from platform
		var (
			osName   = strings.Split(string(platform), "/")[0]
			archName = strings.Split(string(platform), "/")[1]
		)
		binaryName := fmt.Sprintf("bw-%s-%s", osName, archName)
		publisher = publisher.WithExec([]string{
			"gh", "release", "upload", w.ReleaseVersion, binaryName, "--clobber",
		})
	}
	_, err = publisher.Stdout(ctx)
	if err != nil {
		return err
	}

	return nil
}
