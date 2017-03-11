package docker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/andreas-kokkalis/dock-server/pkg/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Repo ...
type Repo struct {
	docker    *DockerCli
	imageRepo string
}

// NewRepo returns a new docker repo
func NewRepo(docker *DockerCli, dockerConfig map[string]string) *Repo {
	return &Repo{docker, dockerConfig["repo"]}
}

/*

	Functionality for Containers

*/

// ContainerCheckState checks if container has the desired state
func (d *Repo) ContainerCheckState(containerID string, state string) (bool, error) {

	var inspect types.ContainerJSON
	var err error

	for i := 0; i < 50; i++ {
		time.Sleep(time.Millisecond)
		inspect, err = d.docker.Cli.ContainerInspect(context.Background(), containerID)
		if err != nil {
			return false, errors.New("Error inspecting state of container: " + containerID)
		}
		fmt.Printf("Container: %s, Status: %s\n", containerID, inspect.State.Status)
		if inspect.State.Status == state {
			return true, nil
		}
	}
	// After X miliseconds container has not started
	return false, nil
}

// ContainerGetUsedPorts returns the list of used ports
func (d *Repo) ContainerGetUsedPorts() (ports map[int]string, err error) {
	// If containers are filtered by status, prepare the ContainerListOptions
	var containerListOptions types.ContainerListOptions

	filterArgs := filters.NewArgs()
	for _, imageRepo := range d.ImageListRepositories() {
		filterArgs.Add("ancestor", imageRepo)
	}
	filterArgs.Add("status", "running")
	containerListOptions = types.ContainerListOptions{Filters: filterArgs}

	containers, err := d.docker.Cli.ContainerList(context.Background(), containerListOptions)
	if err != nil {
		return ports, err
	}

	// Extract containerID, ImageName, and Status
	ports = make(map[int]string)
	for _, container := range containers {
		log.Printf("[PortMapper]: port %v is in use by contaienr %v\n", container.Ports[0].PublicPort, container.ID[:12])
		// containerList[i] = Ctn{ID: container.ID[:10], Image: container.Image, Status: container.Status, State: container.State}
		ports[int(container.Ports[0].PublicPort)] = container.ID[:12]
	}
	return ports, nil
}

// ContainerRemove force removes a container
func (d *Repo) ContainerRemove(containerID string, port int) (err error) {
	// t := time.Duration(time.Millisecond * 100)

	/*
		// First stop the container
		err = Cli.ContainerStop(context.Background(), containerID, &t)
		if err != nil {
			// shut up ...
			fmt.Println("Attemted to stop the container")
			//return err
		}

		// Then kill it
		err = Cli.ContainerKill(context.Background(), containerID, "SIGKILL")
		if err != nil {
			fmt.Println("Attemted to kill the container")
			// srv.FreeResource(srv.PortResources, port)
			// return err
		}
	*/

	// After the container is killed free the port resource
	// XXX: The next line is moved from here.
	// ContainerPorts.Remove(port)
	// rm -f the container
	options := types.ContainerRemoveOptions{Force: true}
	err = d.docker.Cli.ContainerRemove(context.Background(), containerID, options)
	if err != nil {
		log.Printf("[RemoveContainer]: An error occured while removing container %s\n\tError: %v\n", containerID, err.Error())
		return err
	}
	return err
}

/*
// ContainerCreate creates a container based on
// imageName and reference Tag,
// and returns the containerID
func (d *Repo) ContainerCreate(imageID, password string) (containerID string, port int, err error) {
	// Set environment variables for shellinabox container
	envVars := []string{"SIAB_PASSWORD=" + password, "SIAB_SUDO=true"}
	// Get the imageTag
	refTag, err := d.GetTagByID(imageID)
	if err != nil {
		return containerID, port, err
	}

	// --- Container configuration
	// Set container port. This port will be exposed and mapped to a host port
	var natPort nat.Port = "4200/tcp"
	exposedPorts := map[nat.Port]struct{}{natPort: {}}
	// Define configuration required to create a container
	img := conf.GetVal("dc.imagerepo.name") + ":" + refTag
	containerConfig := container.Config{Env: envVars, ExposedPorts: exposedPorts, Image: img}
	// Get a non utilized host port, to avoid collision
	port, err = ContainerPorts.Reserve()
	if err != nil {
		log.Printf("[CreateContainer]: %v", err.Error())
		return "", -1, err
	}
	if port == -1 {
		log.Printf("[CreateContainer]: No ports were available to reserve.\n")
		return "", -1, errors.New("there are no resources available in the system.")
	}

	// --- Host configuration
	// Prepare portBindings containerPort -> Host port. are part of PortMap
	portBindings := []nat.PortBinding{nat.PortBinding{HostPort: strconv.Itoa(port)}}
	// ContainerPorts.PrintUsed() // Debug Logging
	// PortMap is member of container.HostConfig
	portMap := map[nat.Port][]nat.PortBinding{natPort: portBindings}
	hostConfig := container.HostConfig{PortBindings: portMap}

	// Send the request to create the container
	body, err := d.docker.Cli.ContainerCreate(context.Background(), &containerConfig, &hostConfig, &network.NetworkingConfig{}, "")
	if err != nil {
		ContainerPorts.Remove(port)
		log.Printf("[CreateContainer]: %v", err.Error())
		return "", -1, err
	}

	// Return only the first 12 digits from the sha256 identifier of the container
	return body.ID[:12], port, nil
}
*/

/*
	Functionality for Images
*/

// ImageList retrieves the list of docker images from Docker Engine
// via the Docker Remote API.
// It only returns the image ID, the repotags
func (d *Repo) ImageList() ([]api.Img, error) {
	var imageList []api.Img

	// Get list of images from Docker Engine
	// types.ImageListOptions accepts filters.
	// Since no filters are used, all images are returned.
	// XXX: consider limiting this to only the original base image. ==> NOT SUPPORTED

	// the args will be {"image.name":{"ubuntu":true},"label":{"label1=1":true,"label2=2":true}}
	//map[string]map[string]bool

	// From the docker daemon source code:
	/*var acceptedImageFilterTags = map[string]bool{
		"dangling":  true,
		"label":     true,
		"before":    true,
		"since":     true,
		"reference": true,
	}
	*/

	images, err := d.docker.Cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return imageList, err
	}

	// Extract imageID, RepoTags for specific type of images
	for _, image := range images {
		// fmt.Println(image.RepoTags[0])
		s := image.RepoTags[0]
		if s[0:strings.LastIndex(s, ":")] == d.imageRepo {
			log.Printf("[List-Images]: %+v\n", image)
			imageList = append(imageList, api.Img{ID: image.ID[7:19], RepoTags: image.RepoTags, CreatedAt: time.Unix(image.Created, 0)})
		}
	}
	return imageList, nil
}

// ImageListRepositories returns the lsit of the dc repositories and tags
func (d *Repo) ImageListRepositories() (imageList []string) {
	images, err := d.docker.Cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return imageList
	}

	// The repositories
	for _, image := range images {
		// fmt.Printf("%+v\n\n", image.RepoTags)
		s := image.RepoTags[0]
		if s[0:strings.LastIndex(s, ":")] == d.imageRepo {
			imageList = append(imageList, image.RepoTags[0])
		}
	}
	return imageList
}

// ImageTagByID returns an imageTag given and imageID:
// 	performs an ImageList request, gathers all results and returns the tag of the given imageID
func (d *Repo) ImageTagByID(imageID string) (string, error) {
	images, err := d.docker.Cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return "", err
	}
	// Extract imageID, RepoTags for specific type of images
	for _, image := range images {
		if image.ID[7:19] == imageID {
			s := image.RepoTags[0]
			return s[strings.LastIndex(s, ":")+1:], nil
		}
	}
	return "", nil
}

// ImageHistory returns image history
func (d *Repo) ImageHistory(imageID string) ([]api.ImgHistory, error) {
	var history []types.ImageHistory
	var err error

	history, err = d.docker.Cli.ImageHistory(context.Background(), imageID)
	if err != nil {
		return nil, err
	}

	res := []api.ImgHistory{}

	for _, v := range history {
		res = append(res, api.ImgHistory{ID: v.ID[7:19], CreatedAt: time.Unix(v.Created, 0), RepoTags: v.Tags, Comment: v.Comment})
		// Only the first is needed
		break
	}

	return res, nil
}

// RemoveImage removes an image
func (d *Repo) RemoveImage(imageID string) error {

	options := types.ImageRemoveOptions{}
	tag, _ := d.GetTagByID(imageID)
	imgDelete, err := d.docker.Cli.ImageRemove(context.Background(), d.imageRepo+":"+tag, options)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// TODO: figure out what to do with the imgDelete
	fmt.Println(imgDelete)

	return nil
}

// GetTagByID returns an imageTag given and imageID:
// 	performs an ImageList request, gathers all results and returns the tag of the given imageID
func (d *Repo) GetTagByID(imageID string) (string, error) {
	images, err := d.docker.Cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return "", err
	}
	// Extract imageID, RepoTags for specific type of images
	for _, image := range images {
		if image.ID[7:19] == imageID {
			s := image.RepoTags[0]
			return s[strings.LastIndex(s, ":")+1:], nil
		}
	}
	return "", nil
}

// ContainersByImageID returns running containers of specific ImageID
func (d *Repo) ContainersByImageID(imageID string) (containerList []api.Ctn, err error) {
	// If containers are filtered by status, prepare the ContainerListOptions
	var containerListOptions types.ContainerListOptions

	filterArgs := filters.NewArgs()
	filterArgs.Add("ancestor", imageID)
	filterArgs.Add("status", "running")

	containerListOptions = types.ContainerListOptions{Filters: filterArgs}
	containers, err := d.docker.Cli.ContainerList(context.Background(), containerListOptions)
	if err != nil {
		return containerList, err
	}

	// Extract containerID, ImageName, and Status
	containerList = make([]api.Ctn, len(containers))
	for i, container := range containers {
		containerList[i] = api.Ctn{ID: container.ID[:10], Image: container.Image, Status: container.Status, State: container.State}
	}
	return containerList, nil
}
