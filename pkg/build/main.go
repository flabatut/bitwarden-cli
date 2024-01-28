package build

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

// TODO factorize
const containerDistDir = "/dist"

type Workflow struct {
	Client               *dagger.Client
	ReleaseVersion       string
	BuilderNodeJSVersion string
	RunnerEntryPointPath string
	BuilderWorkDir       string
	BuilderImage         string
	RunnerImage          string
	BuilderPlatforms     []dagger.Platform
}

func (w *Workflow) Build(ctx context.Context) ([]*dagger.Container, *dagger.Directory, error) {
	fmt.Println("Building with Dagger")
	var (
		// argonVersion = "0.31.2"
		zipFile             = "cli-" + w.ReleaseVersion + ".zip"
		downloadUrl         = "https://github.com/bitwarden/clients/archive/refs/tags/" + zipFile
		extractedZipDirName = "clients-cli-" + w.ReleaseVersion
	)

	// set the base container
	// set environment variables
	builder := w.Client.Container().
		From(w.BuilderImage).
		WithWorkdir(w.BuilderWorkDir).
		// Checkout release tarball
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "make", "python3", "g++"}).
		WithExec([]string{"wget", "-q", downloadUrl}).
		WithExec([]string{"unzip", "-q", zipFile})

	// Configure global build environment
	builder = builder.
		WithWorkdir(filepath.Join(w.BuilderWorkDir, extractedZipDirName)).
		WithExec([]string{"npm", "ci", "--include", "dev"})

	// Configure cli build env
	builder = builder.
		WithWorkdir(filepath.Join(w.BuilderWorkDir, extractedZipDirName, "/apps/cli")).
		WithExec([]string{"npm", "run", "build:prod"}).
		WithExec([]string{"npm", "run", "clean"})

	// Create one container per cross env
	containerPlatformVariants := make([]*dagger.Container, 0, len(w.BuilderPlatforms))
	for _, platform := range w.BuilderPlatforms {
		// forge npm pkg target platform name
		targetPlatform, err := w.getTargetPlatform(platform)
		if err != nil {
			return nil, nil, err
		}
		// extract os/arch from platform
		var (
			osName   = strings.Split(string(platform), "/")[0]
			archName = strings.Split(string(platform), "/")[1]
		)
		// forge binary output path
		distOutput := filepath.Join(containerDistDir, "bw-"+osName+"-"+archName)
		// cross compile for platform
		builder = builder.WithExec([]string{
			"npx", "pkg", ".", "--public-packages", "--targets", targetPlatform, "--output", distOutput,
		}).WithExec([]string{
			"bash", "-c", fmt.Sprintf("md5sum %s > %s.checksum", distOutput, distOutput),
		})

		// only build docker images for linux supported
		if osName == "linux" {
			// create the image
			runner := w.Client.Container(dagger.ContainerOpts{Platform: platform}).
				From(w.RunnerImage).
				WithFile(w.RunnerEntryPointPath, builder.File(distOutput)).
				WithEntrypoint([]string{w.RunnerEntryPointPath})
			// dont publish yet but append container in returned list
			containerPlatformVariants = append(containerPlatformVariants, runner)
		}
	}

	// keep cross env build dir as artifact
	artifactPlatformVariants := builder.Directory(containerDistDir)
	return containerPlatformVariants, artifactPlatformVariants, nil
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
