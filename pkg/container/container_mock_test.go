package container

import (
	"github.com/docker/docker/api/types/container"
	imageTypes "github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"
	specs "github.com/moby/docker-image-spec/specs-go/v1"
)

type MockContainerUpdate func(*container.InspectResponse, *imageTypes.InspectResponse)

func MockContainer(updates ...MockContainerUpdate) *Container {
	containerInfo := container.InspectResponse{
		ContainerJSONBase: &container.ContainerJSONBase{
			ID:         "container_id",
			Image:      "image",
			Name:       "test-containrrr",
			HostConfig: &container.HostConfig{},
		},
		Config: &container.Config{Labels: map[string]string{}},
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
		portBindings := nat.PortMap{}
		for _, pbs := range portBindingSources {
			portBindings[nat.Port(pbs)] = []nat.PortBinding{}
		}
		c.ContainerJSONBase.HostConfig.PortBindings = portBindings
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
		if c.ContainerJSONBase == nil {
			c.ContainerJSONBase = &container.ContainerJSONBase{}
		}
		if c.ContainerJSONBase.HostConfig == nil {
			c.ContainerJSONBase.HostConfig = &container.HostConfig{}
		}
		c.ContainerJSONBase.HostConfig.Links = links
	}
}

func WithLabels(labels map[string]string) MockContainerUpdate {
	return func(c *container.InspectResponse, i *imageTypes.InspectResponse) {
		c.Config.Labels = labels
	}
}

func WithContainerState(state container.State) MockContainerUpdate {
	return func(cnt *container.InspectResponse, img *imageTypes.InspectResponse) {
		if cnt.ContainerJSONBase == nil {
			cnt.ContainerJSONBase = &container.ContainerJSONBase{}
		}
		cnt.ContainerJSONBase.State = &state
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
