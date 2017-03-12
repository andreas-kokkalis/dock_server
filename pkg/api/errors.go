package api

const (
	// ErrInvalidImageID when imageID is not a sha256
	ErrInvalidImageID = "ImageID is invalid"
	// ErrInvalidContainerID when the id is not a sha256
	ErrInvalidContainerID = "ContainerID is invalid"
	// ErrServerError when an error comes from the backend
	ErrServerError = "Server Error"
	//ErrInvalidContainerState when param is not correct
	ErrInvalidContainerState = "Container state is invalid"
	// InvalidPostData when running a container and post data are insufficient
	InvalidPostData = "POST parameters have errors"
	// DatabaseError occured
	DatabaseError = "A database error occured"

	// ErrUsernameNotExists when username does not exist in db
	ErrUsernameNotExists = "Username does not exist mismatch"
	// ErrPasswordMismatch when given password does not match the stored
	ErrPasswordMismatch = "Password mismatch"
	// CredentialsInvalid When username and password are invalid
	CredentialsInvalid = "Invalid username and/or password"
	// ImageNotFound .
	ImageNotFound = "Image could not be found"
	// ErrContainerAlreadyKilled when the container is not running
	ErrContainerAlreadyKilled = "Container did not exist."
)
