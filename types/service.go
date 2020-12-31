package types

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/fhsinchy/tent/config"

	"github.com/containers/podman/v2/pkg/bindings/containers"
	"github.com/containers/podman/v2/pkg/bindings/images"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/containers/podman/v2/pkg/specgen"
)

// Service describes the properties and methods for a service like MySQL or Redis. All the available services in tent uses this struct as their type.
type Service struct {
	Tag         string
	Name        string
	Image       string
	Volume      specgen.NamedVolume
	PortMapping specgen.PortMapping
	Env         map[string]string
	Command     []string
	HasVolumes  bool
	HasCommand  bool
	Prompts     map[string]bool
}

// CreateContainer method creates a new container with using a given image pulled by PullImage method.
func (service *Service) CreateContainer(connText *context.Context) string {
	var containerID string

	imageExists, err := images.Exists(*connText, service.GetImageName())
	if err != nil {
		log.Fatalln(err)
	}

	if !imageExists {
		_, err := images.Pull(*connText, service.GetImageName(), entities.ImagePullOptions{})
		if err != nil {
			log.Fatalln(err)
		}
	}

	containerExists, err := containers.Exists(*connText, service.GetContainerName(), false)
	if err != nil {
		log.Fatalln(err)
	}

	if !containerExists {
		fmt.Printf("Creating %s container using %s image...\n", service.GetContainerName(), service.GetImageName())
		s := specgen.NewSpecGenerator(service.GetImageName(), false)
		s.Env = service.Env
		s.Remove = true
		s.Name = service.GetContainerName()
		s.PortMappings = append(s.PortMappings, service.PortMapping)

		if service.HasVolumes {
			service.Volume.Name = service.GetVolumeName()
			s.Volumes = append(s.Volumes, &service.Volume)
		}

		if service.HasCommand {
			s.Command = service.Command
		}

		createResponse, err := containers.CreateWithSpec(*connText, s)
		if err != nil {
			log.Fatalln(err)
		}

		containerID = createResponse.ID
	}

	return containerID
}

// ShowPrompt method presents user with user friendly prompts.
func (service *Service) ShowPrompt() {
	if service.Prompts["tag"] {
		var tag string
		fmt.Printf("Which tag you want to use? (default: %s): ", service.Tag)
		fmt.Scanln(&tag)
		if tag != "" {
			service.Tag = tag
		}
	}

	if service.Prompts["port"] {
		var port uint16
		fmt.Printf("Host system port? (default: %d): ", service.PortMapping.HostPort)
		fmt.Scanln(&port)
		if port != 0 {
			service.PortMapping.HostPort = port
		}
	}

	if service.Prompts["username"] {
		var username string
		fmt.Printf("Username for the root user? (default: %s): ", service.Env[config.Envs[service.Name]["username"]])
		fmt.Scanln(&username)
		if username != "" {
			service.Env[config.Envs[service.Name]["username"]] = username
		}
	}

	if service.Prompts["password"] {
		var password string
		fmt.Printf("Password for the root user? (default: %s): ", service.Env[config.Envs[service.Name]["password"]])
		fmt.Scanln(&password)
		if password != "" {
			service.Env[config.Envs[service.Name]["password"]] = password
		}
	}

	if service.Prompts["volume"] {
		var volume string
		fmt.Printf("Volume name for persisting data? (default: %s): ", service.GetVolumeName())
		fmt.Scanln(&volume)
		if volume != "" {
			service.Volume.Name = volume
		}
	}
}

// GetContainerName method generates unique name for each container by combining their image tag and exposed port number.
func (service *Service) GetContainerName() string {
	container := "tent" + "-" + service.Name + "-" + service.Tag + "-" + strconv.Itoa(int(service.PortMapping.HostPort))

	return container
}

// GetVolumeName method generates unique name for each volume used by different containers by using their container name.
func (service *Service) GetVolumeName() string {
	volume := service.GetContainerName() + "-" + "data"

	return volume
}

// GetImageName method generates full image name for services by combining their image name and tag.
func (service *Service) GetImageName() string {
	image := service.Image + ":" + service.Tag

	return image
}
