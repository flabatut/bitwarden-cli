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
		WithEnvVariable("GH_REPO", ghRepo).
		WithSecretVariable("GH_TOKEN", w.RegistryPassword).
		WithExec([]string{"apk", "add", "github-cli"}).
		WithMountedDirectory(containerDistDir, artifacts).
		// WithEnvVariable("GH_DEBUG", "api").
		WithWorkdir(containerDistDir)

	for _, platform := range w.PlatformsVariants {
		// extract os/arch from platform
		var (
			osName   = strings.Split(string(platform), "/")[0]
			archName = strings.Split(string(platform), "/")[1]
		)
		binaryName := fmt.Sprintf("bw-%s-%s", osName, archName)
		md5SumName := fmt.Sprintf("bw-%s-%s.checksum", osName, archName)
		publisher = publisher.WithExec([]string{
			"gh", "release", "upload", w.ReleaseVersion, binaryName, "--clobber",
		}).WithExec([]string{
			"gh", "release", "upload", w.ReleaseVersion, md5SumName, "--clobber",
		})
	}
	_, err := publisher.Stdout(ctx)
	if err != nil {
		return err
	}

	return nil
}
