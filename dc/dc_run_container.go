package dc

import "strconv"

// RunContainer does something
func RunContainer(image, refTag, username, password string) (string, error) {
	id, port, err := CreateContainer(image, refTag, username, password)
	if err != nil {
		return "", err
	}
	err = StartContainer(id)
	if err != nil {
		return "", err
	}

	url := "https://localhost:" + strconv.Itoa(port)
	return url, nil
}
