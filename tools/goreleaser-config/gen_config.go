package main

import (
	"fmt"
	"strings"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

// Generate creates the goreleaser configuration based on the environment.
var validEnvironments = []string{"devspace", "develop", "production"}

func Generate(environment string) config.Project {
	checkEnvironments(environment)

	architectures := []string{"amd64", "arm64"}

	project := config.Project{
		ProjectName: "plugin",
		Version:     2,
		Env:         commonEnv(environment),
		Before: config.Before{
			Hooks: []config.Hook{
				{
					Cmd: "go mod tidy",
				},
				{
					Cmd: "./tools/bin/goreleaser_utils before_hook",
				},
			},
		},
		Builds:          builds(environment),
		Dockers:         dockers(environment, architectures),
		DockerManifests: dockerManifests(environment),
		Checksum: config.Checksum{
			NameTemplate: "checksums.txt",
		},
		Snapshot: config.Snapshot{
			VersionTemplate: "{{ .Env.VERSION }}-{{ .ShortCommit }}",
		},
		Nightly: config.Nightly{
			VersionTemplate: "{{ .Env.VERSION }}-{{ .Env.IMG_TAG }}",
		},
		Partial: config.Partial{
			By: "target",
		},
		Release: config.Release{
			Disable: "true",
		},
		Archives: []config.Archive{
			{
				Format: "binary",
			},
		},
		Changelog: config.Changelog{
			Disable: "true",
		},
	}
	if environment == "devspace" {
		versionTemplate := `v0.0.0-{{ .Runtime.Goarch }}-{{ .Now.Format "2006-01-02-15-04-05Z" }}`
		project.Snapshot = config.Snapshot{VersionTemplate: versionTemplate}
		project.Nightly = config.Nightly{VersionTemplate: versionTemplate}
	}

	// Add SBOMs if needed
	if environment == "production" {
		project.Changelog = config.Changelog{
			Sort: "asc",
			Filters: config.Filters{
				Exclude: []string{
					"^docs:",
					"^test:",
				},
			},
		}
		project.Archives = []config.Archive{
			{
				Format: "tar.gz",
			},
		}
		project.SBOMs = []config.SBOM{
			{
				Artifacts: "archive",
			},
		}
	}

	return project
}

func checkEnvironments(environment string) {
	valid := false
	for _, env := range validEnvironments {
		if environment == env {
			valid = true
			break
		}
	}
	if !valid {
		panic(fmt.Sprintf("invalid environment: %s, valid environments are %v", environment, validEnvironments))
	}
}

// commonEnv returns the common environment variables used across environments.
func commonEnv(environment string) []string {
	envs := []string{
		`IMG_PRE={{ if index .Env "IMAGE_PREFIX"  }}{{ .Env.IMAGE_PREFIX }}{{ else }}localhost:5001{{ end }}`,
		`IMG_TAG={{ if index .Env "IMAGE_TAG" }}{{ .Env.IMAGE_TAG }}{{ else }}develop{{ end }}`,
		`CGO_ENABLED=1`,
	}

	if environment != "devspace" {
		envs = append(envs, `VERSION={{ if index .Env "PLUGIN_VERSION" }}{{ .Env.PLUGIN_VERSION }}{{ else }}v0.0.0-local{{ end }}`)
	}
	return envs
}

// builds returns the build configurations based on the environment.
func builds(environment string) []config.Build {
	switch environment {
	case "devspace":
		return []config.Build{
			build(true),
		}
	case "develop", "production":
		return []config.Build{
			build(false),
		}

	default:
		return nil
	}
}

// build creates a build configuration.
func build(isDevspace bool) config.Build {
	dynamicLinker := `/lib{{ if contains .Runtime.Goarch "amd64" }}64{{end}}/ld-linux-{{ if contains .Runtime.Goarch "arm64" }}aarch64{{ else }}x86-64{{end}}.so.1`

	ldflags := []string{
		"-s -w -r=$ORIGIN/libs",
		"-X github.com/goplugin/pluginv3.0/v2/core/static.Sha={{ .FullCommit }}",
		fmt.Sprintf(`-extldflags "-Wl,--dynamic-linker=%s"`, dynamicLinker),
	}

	if isDevspace {
		ldflags = append(ldflags, "-X github.com/goplugin/pluginv3.0/v2/core/static.Version={{ .Version }}")
	} else {
		ldflags = append(ldflags, "-X github.com/goplugin/pluginv3.0/v2/core/static.Version={{ .Env.VERSION }}")
	}

	return config.Build{
		Binary:          "plugin",
		NoUniqueDistDir: "true",
		Targets:         []string{"go_first_class"},
		Hooks: config.BuildHookConfig{
			Post: []config.Hook{
				{Cmd: "./tools/bin/goreleaser_utils build_post_hook {{ dir .Path }}"},
			},
		},
		BuildDetails: config.BuildDetails{
			Flags:   []string{"-trimpath", "-buildmode=pie"},
			Ldflags: ldflags,
		},
	}
}

// dockers returns the docker configurations based on the environment.
func dockers(environment string, architectures []string) []config.Docker {
	var dockers []config.Docker
	switch environment {
	case "devspace":
		dockers = []config.Docker{
			docker("linux-amd64", "linux", "amd64", environment, true),
			docker("linux-arm64", "linux", "arm64", environment, true),
		}

	case "develop", "production":
		imageNames := []string{"plugin", "ccip"}

		for _, imageName := range imageNames {
			for _, arch := range architectures {
				id := fmt.Sprintf("linux-%s-%s", arch, imageName)
				pluginId := fmt.Sprintf("%s-plugins", id)

				dockers = append(dockers, docker(id, "linux", arch, environment, false))
				dockers = append(dockers, docker(pluginId, "linux", arch, environment, false))
			}
		}
	}
	return dockers
}

// docker creates a docker configuration.
func docker(id, goos, goarch, environment string, isDevspace bool) config.Docker {
	isCCIP := strings.Contains(id, "ccip")
	isPlugins := strings.Contains(id, "plugins")
	extraFiles := []string{"tmp/libs"}
	if isPlugins || isDevspace {
		extraFiles = append(extraFiles, "tmp/plugins")
	}
	if isCCIP {
		extraFiles = append(extraFiles, "ccip/config")
	}

	buildFlagTemplates := []string{
		fmt.Sprintf("--platform=%s/%s", goos, goarch),
		"--pull",
		"--build-arg=PLUGIN_USER=plugin",
		"--build-arg=COMMIT_SHA={{ .FullCommit }}",
	}

	if strings.Contains(id, "ccip") {
		buildFlagTemplates = append(buildFlagTemplates,
			"--build-arg=CL_CHAIN_DEFAULTS=/plugin/ccip-config")
	}

	if strings.Contains(id, "plugins") || isDevspace {
		buildFlagTemplates = append(buildFlagTemplates,
			"--build-arg=CL_MEDIAN_CMD=plugin-feeds",
			"--build-arg=CL_MERCURY_CMD=plugin-mercury",
			"--build-arg=CL_SOLANA_CMD=plugin-solana",
			"--build-arg=CL_STARKNET_CMD=plugin-starknet",
		)
	}

	buildFlagTemplates = append(buildFlagTemplates,
		`--label=org.opencontainers.image.created={{ .Date }}`,
		`--label=org.opencontainers.image.description="node of the decentralized oracle network, bridging on and off-chain computation"`,
		`--label=org.opencontainers.image.licenses=MIT`,
		`--label=org.opencontainers.image.revision={{ .FullCommit }}`,
		`--label=org.opencontainers.image.source=https://github.com/goplugin/pluginv3.0`,
		`--label=org.opencontainers.image.title=plugin`,
		`--label=org.opencontainers.image.url=https://github.com/goplugin/pluginv3.0`,
	)
	if !isDevspace {
		buildFlagTemplates = append(buildFlagTemplates,
			`--label=org.opencontainers.image.version={{ .Env.VERSION }}`,
		)
	}

	dockerConfig := config.Docker{
		ID:                 id,
		Dockerfile:         "core/plugin.goreleaser.Dockerfile",
		Use:                "buildx",
		Goos:               goos,
		Goarch:             goarch,
		Files:              extraFiles,
		BuildFlagTemplates: buildFlagTemplates,
	}

	// We always want to build both versions as a test, but
	// only push the relevant version based on the tag name
	//
	// We also expect the production config file to only be run during a tag push,
	// enforced inside our github actions workflow, "build-publish"
	if environment == "production" {
		if isCCIP {
			dockerConfig.SkipPush = "{{ not (contains .Tag \"-ccip\") }}"
		} else {
			dockerConfig.SkipPush = "{{ contains .Tag \"-ccip\" }}"
		}
	}

	// This section handles the image templates for the docker configuration
	if environment == "devspace" {
		dockerConfig.ImageTemplates = []string{"{{ .Env.IMAGE }}"}
	} else {
		base := "{{ .Env.IMG_PRE }}"
		// On production envs, we have the ECR prefix for the image
		if environment == "production" {
			if isCCIP {
				base = base + "/plugin/plugin-ccip-experimental-goreleaser"
			} else {
				base = base + "/plugin/plugin-experimental-goreleaser"
			}
		} else {
			if isCCIP {
				base = base + "/ccip"
			} else {
				base = base + "/plugin"
			}
		}

		imageTemplates := []string{}
		if strings.Contains(id, "plugins") {
			taggedBase := fmt.Sprintf("%s:{{ .Env.IMG_TAG }}-plugins", base)
			// We have a default, non-arch specific image for plugins that defaults to amd64
			if goarch == "amd64" {
				imageTemplates = append(imageTemplates, taggedBase)
			}
			imageTemplates = append(imageTemplates,
				fmt.Sprintf("%s-%s", taggedBase, archSuffix(id)),
				fmt.Sprintf("%s:sha-{{ .ShortCommit }}-plugins-%s", base, archSuffix(id)))
		} else {
			taggedBase := fmt.Sprintf("%s:{{ .Env.IMG_TAG }}", base)
			// We have a default, non-arch specific image for plugins that defaults to amd64
			if goarch == "amd64" {
				imageTemplates = append(imageTemplates, taggedBase)
			}
			imageTemplates = append(imageTemplates,
				fmt.Sprintf("%s-%s", taggedBase, archSuffix(id)),
				fmt.Sprintf("%s:sha-{{ .ShortCommit }}-%s", base, archSuffix(id)))
		}

		dockerConfig.ImageTemplates = imageTemplates
	}

	return dockerConfig
}

// archSuffix returns the architecture suffix for image tags.
func archSuffix(id string) string {
	if strings.Contains(id, "arm64") {
		return "arm64"
	}
	return "amd64"
}

// dockerManifests returns the docker manifest configurations based on the environment.
func dockerManifests(environment string) []config.DockerManifest {
	if environment == "devspace" {
		return []config.DockerManifest{
			{
				NameTemplate:   "{{ .Env.IMAGE }}",
				ImageTemplates: []string{"{{ .Env.IMAGE }}"},
			},
		}
	}

	// Define the image names based on the environment
	imageNames := []string{"plugin", "ccip"}

	// FIXME: This is duplicated
	if environment == "production" {
		imageNames = []string{"plugin/plugin-experimental-goreleaser", "plugin/plugin-ccip-experimental-goreleaser"}
	}
	var manifests []config.DockerManifest

	for _, imageName := range imageNames {
		fullImageName := fmt.Sprintf("{{ .Env.IMAGE_PREFIX }}/%s", imageName)

		manifestConfigs := []struct {
			ID     string
			Suffix string
		}{
			{ID: "tagged", Suffix: ":{{ .Env.IMG_TAG }}"},
			{ID: "sha", Suffix: ":sha-{{ .ShortCommit }}"},
			{ID: "tagged-plugins", Suffix: ":{{ .Env.IMG_TAG }}-plugins"},
			{ID: "sha-plugins", Suffix: ":sha-{{ .ShortCommit }}-plugins"},
		}
		for _, cfg := range manifestConfigs {
			nameTemplate := fmt.Sprintf("%s%s", fullImageName, cfg.Suffix)
			manifest := config.DockerManifest{
				ID:             strings.ReplaceAll(fmt.Sprintf("%s-%s", cfg.ID, imageName), "/", "-"),
				NameTemplate:   nameTemplate,
				ImageTemplates: manifestImages(nameTemplate),
			}
			if environment == "production" {
				if strings.Contains(nameTemplate, "ccip") {
					manifest.SkipPush = "{{ not (contains .Tag \"-ccip\") }}"
				} else {
					manifest.SkipPush = "{{ contains .Tag \"-ccip\" }}"
				}
			}
			manifests = append(manifests, manifest)
		}
	}

	return manifests
}

// manifestImages generates image templates for docker manifests.
func manifestImages(imageName string) []string {
	architectures := []string{"amd64", "arm64"}
	var images []string
	// Add the default image for tagged images
	if !strings.Contains(imageName, "sha") {
		images = append(images, imageName)
	}
	for _, arch := range architectures {
		images = append(images, fmt.Sprintf("%s-%s", imageName, arch))
	}
	return images
}
