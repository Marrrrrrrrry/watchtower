package container

import (
	specs "github.com/moby/docker-image-spec/specs-go/v1"
	"github.com/moby/moby/api/types/container"
	imageTypes "github.com/moby/moby/api/types/image"
	"github.com/moby/moby/api/types/network"
)

type MockContainerUpdate func(*container.InspectResponse, *imageTypes.InspectResponse)

func MockContainer(updates ...MockContainerUpdate) *Container {
	containerInfo := container.InspectResponse{
		ID:         "container_id",
		Image:      "image",
		Name:       "test-containrrr",
		HostConfig: &container.HostConfig{},
		Config:     &container.Config{Labels: map[string]string{}},
	}
	image := imageTypes.InspectResponse{
		ID:     "image_id",
		Config: &specs.DockerOCIImageConfig{},
	}

	for _, update := range updates {
		update(&containerInfo, &image)
	}
	return NewContainer(&containerInfo, &image)
}

func WithPortBindings(portBindingSources ...string) MockContainerUpdate {
	return func(c *container.InspectResponse, i *imageTypes.InspectResponse) {
		portBindings := network.PortMap{}
		for _, pbs := range portBindingSources {
			portBindings[network.MustParsePort(pbs)] = []network.PortBinding{}
		}
		c.HostConfig.PortBindings = portBindings
	}
}

func WithImageName(name string) MockContainerUpdate {
	return func(c *container.InspectResponse, i *imageTypes.InspectResponse) {
		c.Config.Image = name
		i.RepoTags = append(i.RepoTags, name)
	}
}

func WithLinks(links []string) MockContainerUpdate {
	return func(c *container.InspectResponse, i *imageTypes.InspectResponse) {
		if c.HostConfig == nil {
			c.HostConfig = &container.HostConfig{}
		}
		c.HostConfig.Links = links
	}
}

func WithLabels(labels map[string]string) MockContainerUpdate {
	return func(c *container.InspectResponse, i *imageTypes.InspectResponse) {
		c.Config.Labels = labels
	}
}

func WithContainerState(state container.State) MockContainerUpdate {
	return func(cnt *container.InspectResponse, img *imageTypes.InspectResponse) {
		cnt.State = &state
	}
}

func WithHealthcheck(healthConfig container.HealthConfig) MockContainerUpdate {
	return func(cnt *container.InspectResponse, img *imageTypes.InspectResponse) {
		cnt.Config.Healthcheck = &healthConfig
	}
}

func WithImageHealthcheck(healthConfig container.HealthConfig) MockContainerUpdate {
	return func(cnt *container.InspectResponse, img *imageTypes.InspectResponse) {
		img.Config.Healthcheck = &healthConfig
	}
}
