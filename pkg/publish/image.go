package publish

import (
	"context"
	"fmt"

	// "strings"

	"dagger.io/dagger"
)

func (w *PublishWorkflow) PublishImages(ctx context.Context) error {
	fmt.Println("Publishing with Dagger")
	// var (
	// 	zipFile             = "cli-" + w.ReleaseVersion + ".zip"
	// 	downloadUrl         = "https://github.com/bitwarden/clients/archive/refs/tags/" + zipFile
	// 	extractedZipDirName = "clients-cli-" + w.ReleaseVersion
	// )

	// set the base container
	// set environment variables
	// builder := client.Container().
	// 	From(w.BuilderImage).
	// 	WithEnvVariable("NODE_ENV", "production").
	// 	// Checkout release tarball
	// 	WithWorkdir(w.BuilderWorkDir).
	// 	WithExec([]string{"apt", "install", "-y", "make", "python3", "g++"}).
	// 	WithExec([]string{"wget", downloadUrl}).
	// 	WithExec([]string{"unzip", zipFile}).
	// 	// Configure build environment
	// 	WithWorkdir(filepath.Join(w.BuilderWorkDir, extractedZipDirName)).
	// 	WithExec([]string{"npm", "install", "--include", "dev"}).
	// 	// Build binaries
	// 	WithWorkdir(filepath.Join(w.BuilderWorkDir, extractedZipDirName, "/apps/cli")).
	// 	WithExec([]string{"npm", "run", "build:prod"}).
	// 	WithExec([]string{"npm", "run", "clean"})

	// // platformVariants := make([]*dagger.Container, 0, len(platforms))
	// containerPlatformVariants := make([]*dagger.Container, 0, len(platforms))
	// for _, platform := range platforms {
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
	// 	runner := client.Container(dagger.ContainerOpts{Platform: platform}).
	// 		From(w.RunnerImage).
	// 		WithFile(w.RunnerEntryPointPath, builder.File(distOutput)).
	// 		WithEntrypoint([]string{w.RunnerEntryPointPath})

	// 	// only build docker images for linux supported
	// 	if osName == "linux" {
	// 		containerPlatformVariants = append(containerPlatformVariants, runner)
	// 	}
	// }

	// hack: TOREMOVE
	// containerPlatformVariants := make([]*dagger.Container, 0, len(w.BuilderPlatforms))

	// docker push
	imageDigest, err := w.Client.Container().
		Publish(ctx, w.PublishAddr, dagger.ContainerPublishOpts{
			PlatformVariants: w.PlatformVariants,
		})
	if err != nil {
		return err
	}
	fmt.Println("published multi-platform image with digest", imageDigest)

	return nil
}

// // getTargetPlatform returns the name of a npm pkg compatible target (<nodeVersion>-<os>-<arch>) based
// // on platform name (os/arch)
// func (w *PublishArtifactWorkflow) getTargetPlatform(platform dagger.Platform) (string, error) {
// 	var osName, archName string

// 	switch {
// 	case strings.HasSuffix(string(platform), "arm64"):
// 		archName = "arm64"
// 	case strings.HasSuffix(string(platform), "amd64"):
// 		archName = "x64"
// 	}

// 	switch {
// 	case strings.HasPrefix(string(platform), "darwin"):
// 		osName = "macos"
// 	case strings.HasPrefix(string(platform), "linux"):
// 		osName = "linux"
// 	}

// 	switch {
// 	case osName == "":
// 		return "", fmt.Errorf("os unsupported: %s", platform)
// 	case archName == "":
// 		return "", fmt.Errorf("architecture unsupported: %s", platform)
// 	}
// 	return fmt.Sprintf("%s-%s-%s", w.BuilderNodeJSVersion, osName, archName), nil
// }
