package test

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

type Workflow struct {
	Client               *dagger.Client
	ReleaseVersion       string
	BuilderNodeJSVersion string
	RunnerEntryPointPath string
	BuilderWorkDir       string
	BuilderImage         string
	RunnerImage          string
	BuilderPlatforms     []dagger.Platform
	variantsImages       []*dagger.Container
	variantsArtifacts    []*dagger.Container
}

func (w *Workflow) Test(ctx context.Context) error {
	fmt.Println("Building with Dagger")
	var (
		zipFile             = "cli-" + w.ReleaseVersion + ".zip"
		downloadUrl         = "https://github.com/bitwarden/clients/archive/refs/tags/" + zipFile
		extractedZipDirName = "clients-cli-" + w.ReleaseVersion
	)

	// set the base container
	// set environment variables
	builder := w.Client.Container().
		From(w.BuilderImage).
		WithEnvVariable("NODE_ENV", "production").
		// Checkout release tarball
		WithWorkdir(w.BuilderWorkDir).
		WithExec([]string{"apt", "install", "-y", "make", "python3", "g++"}).
		WithExec([]string{"wget", downloadUrl}).
		WithExec([]string{"unzip", zipFile}).
		// Configure build environment
		WithWorkdir(filepath.Join(w.BuilderWorkDir, extractedZipDirName)).
		WithExec([]string{"npm", "install", "--include", "dev"}).
		// Build binaries
		WithWorkdir(filepath.Join(w.BuilderWorkDir, extractedZipDirName, "/apps/cli")).
		WithExec([]string{"npm", "run", "build:prod"}).
		WithExec([]string{"npm", "run", "clean"})

	// containerPlatformVariants := make([]*dagger.Container, 0, len(w.BuilderPlatforms))
	for _, platform := range w.BuilderPlatforms {
		// forge npm pkg target platform name
		targetPlatform, err := w.getTargetPlatform(platform)
		if err != nil {
			return err
		}
		// extract os/arch from platform
		var (
			osName   = strings.Split(string(platform), "/")[0]
			archName = strings.Split(string(platform), "/")[1]
		)
		// forge binary output path
		distOutput := filepath.Join(w.BuilderWorkDir, "bw-"+osName+"-"+archName)
		// cross compile for platform
		builder = builder.WithExec([]string{"npx", "pkg", ".", "--targets", targetPlatform, "--output", distOutput})

		// create the runner
		runner := w.Client.Container(dagger.ContainerOpts{Platform: platform}).
			From(w.RunnerImage).
			WithFile(w.RunnerEntryPointPath, builder.File(distOutput)).
			WithEntrypoint([]string{w.RunnerEntryPointPath})

		// only build docker images for linux supported
		if osName == "linux" {
			w.variantsImages = append(w.variantsImages, runner)
		}
		// all os/arch supported for binaries artifacts
		w.variantsArtifacts = append(w.variantsArtifacts, runner)
	}
	return nil
}

// getTargetPlatform returns the name of a npm pkg compatible target (<nodeVersion>-<os>-<arch>) based
// on platform name (os/arch)
func (w *Workflow) getTargetPlatform(platform dagger.Platform) (string, error) {
	var osName, archName string

	switch {
	case strings.HasSuffix(string(platform), "arm64"):
		archName = "arm64"
	case strings.HasSuffix(string(platform), "amd64"):
		archName = "x64"
	}

	switch {
	case strings.HasPrefix(string(platform), "darwin"):
		osName = "macos"
	case strings.HasPrefix(string(platform), "linux"):
		osName = "linux"
	}

	switch {
	case osName == "":
		return "", fmt.Errorf("os unsupported: %s", platform)
	case archName == "":
		return "", fmt.Errorf("architecture unsupported: %s", platform)
	}
	return fmt.Sprintf("%s-%s-%s", w.BuilderNodeJSVersion, osName, archName), nil
}

func (w *Workflow) GetImageVariants() []*dagger.Container {
	return w.variantsImages
}

func (w *Workflow) GetArtifactVariants() []*dagger.Container {
	return w.variantsArtifacts
}
