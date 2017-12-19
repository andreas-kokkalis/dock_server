package integration

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const topDir = "../../"

func TestCompareRegexJSON(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	workDir, _ := os.Getwd()
	testFiles := path.Join(workDir, "testfiles")

	type testCase struct {
		name, file, json, dir, resp string
		expectErr                   bool
	}
	mutations := []testCase{
		testCase{
			name:      "File exists, diff correct",
			file:      path.Join(testFiles, "api_test_good.json"),
			json:      `{"foo":"bar"}`,
			dir:       topDir,
			expectErr: false,
		},
		testCase{
			name:      "File does not exist",
			file:      path.Join(testFiles, "does_not_exist.json"),
			json:      `{"foo":"bar"}`,
			dir:       topDir,
			expectErr: true,
		},
		testCase{
			name:      "empty file",
			file:      path.Join(testFiles, "empty_file.json"),
			json:      `{"foo":"bar"}`,
			dir:       topDir,
			expectErr: true,
		},
		testCase{
			name:      "cannot find python script",
			file:      path.Join(testFiles, "wrong_diff.json"),
			json:      `{"foo":"bar"}`,
			dir:       ".",
			expectErr: true,
		}, testCase{
			name:      "json does not match",
			file:      path.Join(testFiles, "wrong_diff.json"),
			json:      `{"foo":"bar"}`,
			dir:       topDir,
			resp:      "jsondiff - INFO - Changed: foo to bar from wrong",
			expectErr: false,
		},
	}

	for _, test := range mutations {
		t.Run(test.name, func(t *testing.T) {
			log, err := CompareRegexJSON(test.file, test.json, test.dir)
			if test.expectErr {
				assert.Error(err)
			} else {
				assert.Nil(err)
				if test.resp != "" {
					assert.Contains(log, test.resp)
				} else {
					assert.Empty(log)
				}
			}
		})

		for _, test := range mutations {
			t.Run(test.name, func(t *testing.T) {
				log, err := CompareJSON(test.file, test.json, test.dir)
				if test.expectErr {
					assert.Error(err)
				} else {
					assert.Nil(err)
					if test.resp != "" {
						assert.Contains(log, test.resp)
					} else {
						assert.Empty(log)
					}
				}
			})
		}
	}
}
