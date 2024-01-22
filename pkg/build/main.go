package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

var platforms = []dagger.Platform{
	"darwin/amd64", // a.k.a. x86_64
	"darwin/arm64", // a.k.a. aarch64
	"linux/amd64",  // a.k.a. x86_64
	"linux/arm64",  // a.k.a. aarch64
}

const (
	workDir     = "/build"
	publishAddr = "ghcr.io/flabatut/bitwarden-cli:latest"
	// publishAddr         = "localhost:5000/toto:latest"
	releaseVersion      = "v2024.1.0"
	zipFile             = "cli-" + releaseVersion + ".zip"
	downloadUrl         = "https://github.com/bitwarden/clients/archive/refs/tags/" + zipFile
	builderImage        = "mcr.microsoft.com/devcontainers/typescript-node:1-20-bullseye"
	extractedZipDirName = "clients-cli-" + releaseVersion
	runnerImage         = "docker.io/debian:bullseye-slim"
	entrypoint          = "/entrypoint"
	defaultNodeVersion  = "latest"
)

func Build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// set the base container
	// set environment variables
	builder := client.Container().
		From(builderImage).
		WithEnvVariable("NODE_ENV", "production").
		// Checkout release tarball
		WithWorkdir(workDir).
		WithExec([]string{"apt", "install", "-y", "make", "python3", "g++"}).
		WithExec([]string{"wget", downloadUrl}).
		WithExec([]string{"unzip", zipFile}).
		// Configure build environment
		WithWorkdir(filepath.Join(workDir, extractedZipDirName)).
		WithExec([]string{"npm", "install", "--include", "dev"}).
		// Build binaries
		WithWorkdir(filepath.Join(workDir, extractedZipDirName, "/apps/cli")).
		WithExec([]string{"npm", "run", "build:prod"}).
		WithExec([]string{"npm", "run", "clean"})

	// platformVariants := make([]*dagger.Container, 0, len(platforms))
	containerPlatformVariants := make([]*dagger.Container, 0, len(platforms))
	for _, platform := range platforms {
		// forge npm pkg target platform name
		targetPlatform, err := getTargetPlatform(platform)
		if err != nil {
			panic(err)
		}
		// extract os/arch from platform
		var (
			osName   = strings.Split(string(platform), "/")[0]
			archName = strings.Split(string(platform), "/")[1]
		)
		// forge binary output path
		distOutput := filepath.Join(workDir, "bw-"+osName+"-"+archName)
		// cross compile for platform
		builder = builder.WithExec([]string{"npx", "pkg", ".", "--targets", targetPlatform, "--output", distOutput})

		// create the runner
		runner := client.Container(dagger.ContainerOpts{Platform: platform}).
			From(runnerImage).
			WithFile(entrypoint, builder.File(distOutput)).
			WithEntrypoint([]string{entrypoint})

		// only build docker images for linux supported
		if osName == "linux" {
			containerPlatformVariants = append(containerPlatformVariants, runner)
		}
	}

	// docker push
	imageDigest, err := client.Container().
		Publish(ctx, publishAddr, dagger.ContainerPublishOpts{
			PlatformVariants: containerPlatformVariants,
		})
	if err != nil {
		panic(err)
	}
	fmt.Println("published multi-platform image with digest", imageDigest)

	return nil
}

// getTargetPlatform returns the name of a npm pkg compatible target (<nodeVersion>-<os>-<arch>) based
// on platform name (os/arch)
func getTargetPlatform(platform dagger.Platform) (string, error) {
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
	return fmt.Sprintf("%s-%s-%s", defaultNodeVersion, osName, archName), nil
}
