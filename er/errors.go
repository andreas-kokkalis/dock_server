package er

const (
	// InvalidImageID when imageID is not a sha256
	InvalidImageID = "ImageID is invalid"
	// InvalidContainerID when the id is not a sha256
	InvalidContainerID = "ContainerID is invalid"
	// ServerError when an error comes from the backend
	ServerError = "Server Error"
	//InvalidContainerState when param is not correct
	InvalidContainerState = "Container state is invalid"
	// InvalidPostData when running a container and post data are insufficient
	InvalidPostData = "POST parameters have errors"
)
