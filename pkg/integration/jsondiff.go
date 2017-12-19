package integration

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// jsondiff.go Utilizes the command line tool jsondiff.py
// https://github.com/bazaarvoice/json-regex-difftool

const (
	diffExecutable = "pkg/integration/jsondiff.py"
)

var (
	// ErrFileNotFound indicates that a file cannot be found at a given path
	ErrFileNotFound = errors.New("Cannot open file")
)

// CompareRegexJSON invokes the jsondiff.py command line tool to compare json files that contain regex
// expected parameter is the filePath of the expected JSON model.
// actual parameter is the string output as it is usually returned by the API
// topDir is the location of the top level directory, relative to the location of the Spec file that invoked the function
// uses the --use_model flag to evaluate regex contained in the expected file
// Use this function only when you need to compare against a model that contains regular expressions for attribute values.
func CompareRegexJSON(expected string, actual string, topDir string) (string, error) {

	// Check whether the expected file exists
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		return "", errors.New("Cannot find expected JSON model file: " + expected)
	}

	// Create the actual file
	f, err := ioutil.TempFile(os.TempDir(), "tmp_integration")
	if err != nil {
		return "", err
	}

	if _, err = f.WriteString(actual); err != nil {
		return "", err
	}

	if err = f.Close(); err != nil {
		return "", err
	}

	defer func() {
		err = os.Remove(f.Name())
	}()

	script := topDir + diffExecutable
	compareCmd := exec.Command(
		"python",
		script,
		"--use_model",
		f.Name(),
		expected)
	out, err := compareCmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	stringOut := string(out)
	if strings.Contains(stringOut, "jsondiff - INFO - True") {
		return "", nil
	} else if strings.Contains(stringOut, "jsondiff - INFO - False") {
		// does not match compute diff and return
		diffCmd := exec.Command(
			"python",
			script,
			"--diff",
			"--use_model",
			f.Name(),
			expected)
		out, err = diffCmd.CombinedOutput()
		if err != nil {
			return "", err
		}
		return string(out), nil
	} else {
		return "", errors.New("jsondiff returned unexpected output in compare mode")
	}
}

// CompareJSON invokes the jsondiff.py command line tool to compare json files
// expected parameter is the filePath of the expected JSON model.
// actual parameter is the string output as it is usually returned by the API
// topDir is the location of the top level directory, relative to the location of the Spec file that invoked the function
// Use this function only when you need to perform a plain diff of JSON objects.
func CompareJSON(expected string, actual string, topDir string) (string, error) {

	// Check whether the expected file exists
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		return "", errors.New("Cannot find expected JSON model file: " + expected)
	}

	f, err := ioutil.TempFile(os.TempDir(), "tmp_integration")
	if err != nil {
		return "", err
	}

	if _, err = f.WriteString(actual); err != nil {
		return "", err
	}

	if err = f.Close(); err != nil {
		return "", err
	}

	defer func() {
		err = os.Remove(f.Name())
	}()

	script := topDir + diffExecutable
	compareCmd := exec.Command(
		"python",
		script,
		f.Name(),
		expected)
	out, err := compareCmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	stringOut := string(out)
	if strings.Contains(stringOut, "jsondiff - INFO - True") {
		return "", nil
	} else if strings.Contains(stringOut, "jsondiff - INFO - False") {
		// does not match compute diff and return
		diffCmd := exec.Command(
			"python",
			script,
			"--diff",
			f.Name(),
			expected)
		out, err = diffCmd.CombinedOutput()
		if err != nil {
			return "", err
		}
		return string(out), nil
	} else {
		return "", errors.New("jsondiff returned unexpected output in compare mode")
	}
}
