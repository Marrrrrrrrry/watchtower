package mocks

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Marrrrrrrrry/watchtower/pkg/container"
	wt "github.com/Marrrrrrrrry/watchtower/pkg/types"
	dockerContainer "github.com/moby/moby/api/types/container"
	imageTypes "github.com/moby/moby/api/types/image"
	"github.com/moby/moby/api/types/network"
)

func CreateMockContainer(id string, name string, image string, created time.Time) wt.Container {
	content := dockerContainer.InspectResponse{
		ID:      id,
		Image:   image,
		Name:    name,
		Created: created.String(),
		HostConfig: &dockerContainer.HostConfig{
			PortBindings: network.PortMap{},
		},
		Config: &dockerContainer.Config{
			Image:        image,
			Labels:       make(map[string]string),
			ExposedPorts: network.PortSet{},
		},
	}
	return container.NewContainer(
		&content,
		CreateMockImageInfo(image),
	)
}

func CreateMockImageInfo(image string) *imageTypes.InspectResponse {
	return &imageTypes.InspectResponse{
		ID: image,
		RepoDigests: []string{
			image,
		},
	}
}

func CreateMockContainerWithImageInfo(id string, name string, image string, created time.Time, imageInfo imageTypes.InspectResponse) wt.Container {
	return CreateMockContainerWithImageInfoP(id, name, image, created, &imageInfo)
}

func CreateMockContainerWithImageInfoP(id string, name string, image string, created time.Time, imageInfo *imageTypes.InspectResponse) wt.Container {
	content := dockerContainer.InspectResponse{
		ID:         id,
		Image:      image,
		Name:       name,
		Created:    created.String(),
		HostConfig: &dockerContainer.HostConfig{},
		Config: &dockerContainer.Config{
			Image:  image,
			Labels: make(map[string]string),
		},
	}
	return container.NewContainer(
		&content,
		imageInfo,
	)
}

func CreateMockContainerWithDigest(id string, name string, image string, created time.Time, digest string) wt.Container {
	c := CreateMockContainer(id, name, image, created)
	c.ImageInfo().RepoDigests = []string{digest}
	return c
}

func CreateMockContainerWithConfig(id string, name string, image string, running bool, restarting bool, created time.Time, config *dockerContainer.Config) wt.Container {
	content := dockerContainer.InspectResponse{
		ID:    id,
		Image: image,
		Name:  name,
		State: &dockerContainer.State{
			Running:    running,
			Restarting: restarting,
		},
		Created:    created.String(),
		HostConfig: &dockerContainer.HostConfig{PortBindings: network.PortMap{}},
		Config:     config,
	}
	return container.NewContainer(
		&content,
		CreateMockImageInfo(image),
	)
}

func CreateContainerForProgress(index int, idPrefix int, nameFormat string) (wt.Container, wt.ImageID) {
	indexStr := strconv.Itoa(idPrefix + index)
	mockID := indexStr + strings.Repeat("0", 61-len(indexStr))
	contID := "c79" + mockID
	contName := fmt.Sprintf(nameFormat, index+1)
	oldImgID := "01d" + mockID
	newImgID := "d0a" + mockID
	imageName := fmt.Sprintf("mock/%s:latest", contName)
	config := &dockerContainer.Config{
		Image: imageName,
	}
	c := CreateMockContainerWithConfig(contID, contName, oldImgID, true, false, time.Now(), config)
	return c, wt.ImageID(newImgID)
}

func CreateMockContainerWithLinks(id string, name string, image string, created time.Time, links []string, imageInfo *imageTypes.InspectResponse) wt.Container {
	content := dockerContainer.InspectResponse{
		ID:      id,
		Image:   image,
		Name:    name,
		Created: created.String(),
		Config: &dockerContainer.Config{
			Image:  image,
			Labels: make(map[string]string),
		},
	}
	return container.NewContainer(
		&content,
		imageInfo,
	)
}
