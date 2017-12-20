package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const topDir = "../../"

func TestCompareRegexJSON(t *testing.T) {
	tests := []struct {
		name, file, json, dir, resp string
		expectErr                   bool
	}{
		{
			name:      "File exists, diff correct",
			file:      `{"foo":"bar"}`,
			json:      `{"foo":"bar"}`,
			dir:       topDir,
			expectErr: false,
		},
		{
			name:      "empty file",
			file:      "",
			json:      `{"foo":"bar"}`,
			dir:       topDir,
			expectErr: true,
		},
		{
			name:      "cannot find python script",
			file:      `{"foo":"bar"}`,
			json:      `{"foo":"bar"}`,
			dir:       ".",
			expectErr: true,
		},
		{
			name:      "json does not match",
			file:      `{"foo":"wrong"}`,
			json:      `{"foo":"bar"}`,
			dir:       topDir,
			resp:      "jsondiff - INFO - Changed: foo to bar from wrong",
			expectErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log, err := CompareRegexJSON(test.file, test.json, test.dir)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				if test.resp != "" {
					assert.Contains(t, log, test.resp)
				} else {
					assert.Empty(t, log)
				}
			}
		})

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				log, err := CompareJSON(test.file, test.json, test.dir)
				if test.expectErr {
					assert.Error(t, err)
				} else {
					assert.Nil(t, err)
					if test.resp != "" {
						assert.Contains(t, log, test.resp)
					} else {
						assert.Empty(t, log)
					}
				}
			})
		}
	}
}
