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
	// DatabaseError occured
	DatabaseError = "A database error occured"

	// UsernameNotExists when username does not exist in db
	UsernameNotExists = "Username does not exist mismatch"
	// PasswordMismatch when given password does not match the stored
	PasswordMismatch = "Password mismatch"
	// CredentialsInvalid When username and password are invalid
	CredentialsInvalid = "Invalid username and/or password"
	// ImageNotFound .
	ImageNotFound = "Image could not be found"
	// ContainerAlreadyKilled when the container is not running
	ContainerAlreadyKilled = "Container did not exist."
)
