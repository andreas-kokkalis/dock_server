package integration

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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

// writeTempFile accepts a string content and writes it to a temporary file.
// It is the callers responsibility to remove the file from the os temp directory.
func writeTempFile(filePrefix string, content string) (f *os.File, err error) {
	if f, err = ioutil.TempFile(os.TempDir(), filePrefix); err != nil {
		return nil, err
	}
	if _, err = f.WriteString(content); err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	return f, nil
}

// CompareRegexJSON invokes the jsondiff.py command line tool to compare json files that contain regex
// expected parameter is the filePath of the expected JSON model.
// actual parameter is the string output as it is usually returned by the API
// topDir is the location of the top level directory, relative to the location of the Spec file that invoked the function
// uses the --use_model flag to evaluate regex contained in the expected file
// Use this function only when you need to compare against a model that contains regular expressions for attribute values.
func CompareRegexJSON(expected string, actual string, topDir string) (string, error) {

	var fActual, fExpected *os.File
	var err error
	if fActual, err = writeTempFile("tmp_actual", actual); err != nil {
		return "", err
	}
	if fExpected, err = writeTempFile("tmp_expected", expected); err != nil {
		return "", err
	}
	defer func() {
		_ = os.Remove(fActual.Name())
		_ = os.Remove(fExpected.Name())
	}()

	script := path.Join(topDir, diffExecutable)
	compareCmd := exec.Command(
		"python",
		script,
		"--use_model",
		fActual.Name(),
		fExpected.Name())
	diffCmd := exec.Command(
		"python",
		script,
		"--diff",
		"--use_model",
		fActual.Name(),
		fExpected.Name())
	return runJSONDiff(compareCmd, diffCmd)
}

func runJSONDiff(compareCmd *exec.Cmd, diffCmd *exec.Cmd) (string, error) {
	out, err := compareCmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	stringOut := string(out)
	if strings.Contains(stringOut, "jsondiff - INFO - True") {
		return "", nil
	} else if strings.Contains(stringOut, "jsondiff - INFO - False") {
		// does not match compute diff and return
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

	var fActual, fExpected *os.File
	var err error
	if fActual, err = writeTempFile("tmp_actual", actual); err != nil {
		return "", err
	}
	if fExpected, err = writeTempFile("tmp_expected", expected); err != nil {
		return "", err
	}
	defer func() {
		_ = os.Remove(fActual.Name())
		_ = os.Remove(fExpected.Name())
	}()

	script := path.Join(topDir, diffExecutable)
	compareCmd := exec.Command(
		"python",
		script,
		fActual.Name(),
		fExpected.Name())
	diffCmd := exec.Command(
		"python",
		script,
		"--diff",
		fActual.Name(),
		fExpected.Name())
	return runJSONDiff(compareCmd, diffCmd)
}
