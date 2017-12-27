package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// DockerRepo ...
type DockerRepo struct {
	docker    *docker.APIClient
	imageRepo string
}

// NewDockerRepository returns a new docker repo
func NewDockerRepository(docker *docker.APIClient, dockerConfig map[string]string) DockerRepository {
	return &DockerRepo{docker, dockerConfig["repo"]}
}

/*

	Functionality for Containers

*/

// containerCheckState checks if container has the desired state
func (d *DockerRepo) containerCheckState(containerID string, state string) (bool, error) {

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
func (d *DockerRepo) ContainerGetUsedPorts() (ports map[int]string, err error) {
	// If containers are filtered by status, prepare the ContainerListOptions
	var containerListOptions types.ContainerListOptions

	filterArgs := filters.NewArgs()
	for _, imageRepo := range d.imageListRepositories() {
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
func (d *DockerRepo) ContainerRemove(containerID string, port int) (err error) {
	// t := time.Duration(time.Millisecond * 100)

	/*
		// First stop the container
		err = Cli.ContainerStop(context.Background(), containerID, &t)
		if err != nil {
			// shut up ...
			fmt.Println("Attempted to stop the container")
			//return err
		}

		// Then kill it
		err = Cli.ContainerKill(context.Background(), containerID, "SIGKILL")
		if err != nil {
			fmt.Println("Attempted to kill the container")
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
		log.Printf("[RemoveContainer]: An error occurred while removing container %s\n\tError: %v\n", containerID, err.Error())
		return err
	}
	return err
}

// ContainerRun does something
func (d *DockerRepo) ContainerRun(imageID, username, password string, port int) (cfg api.RunConfig, err error) {
	// Create the container
	id /*, port*/, err := d.containerCreate(imageID, password, port)
	if err != nil {
		log.Printf("[RunContainer]: Error while creating: %v\n", err.Error())
		return cfg, err
	}
	// Start the container
	err = d.containerStart(id)
	if err != nil {
		log.Printf("[RunContainer]: Error while starting: %v\n", err.Error())
		return cfg, err
	}

	cfg = api.RunConfig{
		ContainerID: id,
		Username:    username,
		Password:    password,
		Port:        strconv.Itoa(port),
		URL:         "https://127.0.0.1:" + strconv.Itoa(port),
	}
	return cfg, nil
}

// containerCreate creates a container based on
// imageName and reference Tag,
// and returns the containerID
func (d *DockerRepo) containerCreate(imageID, password string, port int) (containerID string /*, port int*/, err error) {
	// Set environment variables for shellinabox container
	envVars := []string{"SIAB_PASSWORD=" + password, "SIAB_SUDO=true"}
	// Get the imageTag
	refTag, err := d.ImageGetTagByID(imageID)
	if err != nil {
		return containerID /*, port*/, err
	}

	// --- Container configuration
	// Set container port. This port will be exposed and mapped to a host port
	var natPort nat.Port = "4200/tcp"
	exposedPorts := map[nat.Port]struct{}{natPort: {}}
	// Define configuration required to create a container
	img := d.imageRepo + ":" + refTag
	containerConfig := container.Config{Env: envVars, ExposedPorts: exposedPorts, Image: img}

	/*
		// Get a non utilized host port, to avoid collision
		port, err = ContainerPorts.Reserve()
		if err != nil {
			log.Printf("[CreateContainer]: %v", err.Error())
			return "", -1, err
		}
		if port == -1 {
			log.Printf("[CreateContainer]: No ports were available to reserve.\n")
			return "", -1, errors.New("there are no resources available in the system")
		}
	*/

	// --- Host configuration
	// Prepare portBindings containerPort -> Host port. are part of PortMap
	portBindings := []nat.PortBinding{{HostPort: strconv.Itoa(port)}}
	// ContainerPorts.PrintUsed() // Debug Logging
	// PortMap is member of container.HostConfig
	portMap := map[nat.Port][]nat.PortBinding{natPort: portBindings}
	hostConfig := container.HostConfig{PortBindings: portMap}

	// Send the request to create the container
	body, err := d.docker.Cli.ContainerCreate(context.Background(), &containerConfig, &hostConfig, &network.NetworkingConfig{}, "")
	if err != nil {
		// ContainerPorts.Remove(port)
		log.Printf("[CreateContainer]: %v", err.Error())
		return "" /*, -1*/, err
	}

	// Return only the first 12 digits from the sha256 identifier of the container
	return body.ID[:12] /*, port*/, nil
}

// containerStart sends a request to start a container
func (d *DockerRepo) containerStart(containerID string) error {

	// Start container
	err := d.docker.Cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		// No need to validate if port number is -1, since error is returned in that case
		return err
	}

	// Check if container is running
	var isRunning bool
	isRunning, err = d.containerCheckState(containerID, "running")
	if err != nil {
		return err
	}
	if !isRunning {
		return errors.New("container not started")
	}
	return nil
}

// ContainerCommit creates a new image from a running container
func (d *DockerRepo) ContainerCommit(comment, author, containerID, refTag string) (string, error) {
	// XXX: new version of docker-daemon adds the docker.io/ domain by default when tagging images.
	fullRef := "docker.io/" + d.imageRepo + ":" + refTag

	// TODO: on options, can add a slice of string with the list of changes for this commit
	options := types.ContainerCommitOptions{
		Comment:   comment,
		Author:    author,
		Reference: fullRef,
	}
	log.Printf("options: %+v", options)
	response, err := d.docker.Cli.ContainerCommit(context.Background(), containerID, options)
	if err != nil {
		fmt.Printf("response: %+v\nerror: %s", response, err.Error())
		return "", err
	}
	fmt.Printf("%+v\n", response)
	log.Printf("[CommitContainer]: Committed container with ID:%s\n", containerID)
	//sha256:baa8ace946df92b5fb1722538d73531503485535604863e34e174a5d284a601b

	imgID := response.ID[7:19]
	return imgID, nil
}

// ContainerList returns the list of containers. Use
// type.ContainerListOptions to filter for state such as
// status=(created,	restarting, running, paused, exited, dead)
// TODO: filter containers by base repository
func (d *DockerRepo) ContainerList(status string) ([]api.Ctn, error) {
	var containerList []api.Ctn

	// If containers are filtered by status, prepare the ContainerListOptions
	var containerListOptions types.ContainerListOptions
	if status != "" {
		filterArgs := filters.NewArgs()
		filterArgs.Add("status", status)
		containerListOptions = types.ContainerListOptions{Filters: filterArgs}

	}
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

/*
	Functionality for Images
*/

// ImageList retrieves the list of docker images from Docker Engine
// via the Docker Remote API.
// It only returns the image ID, the repotags
func (d *DockerRepo) ImageList() ([]api.Img, error) {
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
			// log.Printf("[List-Images]: %+v\n", image)
			imageList = append(imageList, api.Img{ID: image.ID[7:19], RepoTags: image.RepoTags, CreatedAt: time.Unix(image.Created, 0)})
		}
	}
	return imageList, nil
}

// imageListRepositories returns the lsit of the dc repositories and tags
func (d *DockerRepo) imageListRepositories() (imageList []string) {
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

// imageTagByID returns an imageTag given and imageID:
// 	performs an ImageList request, gathers all results and returns the tag of the given imageID
func (d *DockerRepo) imageTagByID(imageID string) (string, error) {
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
func (d *DockerRepo) ImageHistory(imageID string) ([]api.ImgHistory, error) {
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

// ImageRemove removes an image
func (d *DockerRepo) ImageRemove(imageID string) error {

	options := types.ImageRemoveOptions{}
	tag, _ := d.ImageGetTagByID(imageID)
	imgDelete, err := d.docker.Cli.ImageRemove(context.Background(), d.imageRepo+":"+tag, options)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// TODO: figure out what to do with the imgDelete
	fmt.Println(imgDelete)

	return nil
}

// ImageGetTagByID returns an imageTag given and imageID:
// 	performs an ImageList request, gathers all results and returns the tag of the given imageID
func (d *DockerRepo) ImageGetTagByID(imageID string) (string, error) {
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

// GetRunningContainersByImageID returns running containers of specific ImageID
func (d *DockerRepo) GetRunningContainersByImageID(imageID string) (containerList []api.Ctn, err error) {
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
