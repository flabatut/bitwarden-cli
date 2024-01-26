package lint

import (
	"context"
	"fmt"

	// "path/filepath"
	// "strings"

	// "os"

	"dagger.io/dagger"
)

type Workflow struct {
	Client        *dagger.Client
	LinterImage   string
	LinterWorkDir string

	// ReleaseVersion       string
	// BuilderNodeJSVersion string
	// RunnerEntryPointPath string
	// RunnerImage          string
	// BuilderPlatforms     []dagger.Platform
	// RegistryFQDN         string
	// ProjectNamespace     string
	// RegistryUsername     string
	// RegistryPassword     *dagger.Secret
}

func (w *Workflow) Lint(ctx context.Context) error {

	// get build context directory
	contextDir := w.Client.Host().Directory(".")

	// build using Dockerfile
	// publish the resulting container to a registry
	ctr, err := w.Client.
		Container().
		From(w.LinterImage).
		WithDirectory("/build", contextDir).
		WithWorkdir("/build").
		WithExec([]string{"golangci-lint", "-v", "run", "./..."}).
		Stdout(ctx)
	// WithUnixSocket("/var/run/docker.sock", w.Client.Host().UnixSocket("/var/run/docker.sock")).
	// Publish(ctx, "test:latest")
	if err != nil {
		return err
	}
	fmt.Println("Published image to", ctr)

	// get reference to the local project
	// projectDir := w.Client.Host().Directory(".")
	// set the base container
	// set environment variables
	// linter := w.Client.Container().
	// 	From(w.LinterImage).
	// Checkout release tarball
	// WithDirectory(".", projectDir).
	// WithWorkdir(w.LinterWorkDir)
	// WithExec([]string{"apt", "install", "-y", "make", "python3", "g++"}).
	// WithExec([]string{"wget", downloadUrl}).
	// WithExec([]string{"unzip", zipFile}).
	// // Configure build environment
	// WithWorkdir(filepath.Join(w.BuilderWorkDir, extractedZipDirName)).
	// WithExec([]string{"npm", "install", "--include", "dev"}).
	// // Build binaries
	// WithWorkdir(filepath.Join(w.BuilderWorkDir, extractedZipDirName, "/apps/cli")).
	// WithExec([]string{"npm", "run", "build:prod"}).
	// WithExec([]string{"npm", "run", "clean"})

	// containerPlatformVariants := make([]*dagger.Container, 0, len(w.BuilderPlatforms))
	// for _, platform := range w.BuilderPlatforms {
	// 	// forge npm pkg target platform name
	// 	targetPlatform, err := w.getTargetPlatform(platform)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// extract os/arch from platform
	// 	var (
	// 		osName   = strings.Split(string(platform), "/")[0]
	// 		archName = strings.Split(string(platform), "/")[1]
	// 	)
	// 	// forge binary output path
	// 	distOutput := filepath.Join(w.BuilderWorkDir, "bw-"+osName+"-"+archName)
	// 	// cross compile for platform
	// 	builder = builder.WithExec([]string{"npx", "pkg", ".", "--targets", targetPlatform, "--output", distOutput})

	// 	// create the runner
	// 	runner := w.Client.Container(dagger.ContainerOpts{Platform: platform}).
	// 		From(w.RunnerImage).
	// 		WithFile(w.RunnerEntryPointPath, builder.File(distOutput)).
	// 		WithEntrypoint([]string{w.RunnerEntryPointPath})

	// 	// only build docker images for linux supported
	// 	if osName == "linux" {
	// 		containerPlatformVariants = append(containerPlatformVariants, runner)
	// 	}
	// }

	// // docker push
	// // TODO: validate registry URL
	// // TODO: support image name option
	// publishAddress := fmt.Sprintf("%s/%s:%s", w.RegistryFQDN, w.ProjectNamespace, w.ReleaseVersion)
	// imageDigest, err := w.Client.Container().
	// 	WithRegistryAuth(w.RegistryFQDN, w.RegistryUsername, w.RegistryPassword).
	// 	Publish(ctx, publishAddress, dagger.ContainerPublishOpts{
	// 		PlatformVariants: containerPlatformVariants,
	// 	})
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("published multi-platform image with digest", imageDigest)

	return nil
}
