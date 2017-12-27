package image

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories/repomocks"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestListImages(t *testing.T) {
	imgList := []api.Img{
		api.Img{
			ID:       "12345",
			RepoTags: []string{"ref"},
		},
	}
	tests := []struct {
		service        Service
		expectCode     int
		expectResponse interface{}
		name           string
	}{
		{
			service: NewService(nil,
				repomocks.NewDockerRepositoryMock().
					WithImageList(imgList, nil)),
			expectCode:     http.StatusOK,
			expectResponse: imgList,
			name:           "list images",
		},
		{
			service: NewService(nil,
				repomocks.NewDockerRepositoryMock().
					WithImageList([]api.Img{}, errors.New("error"))),
			expectCode:     http.StatusInternalServerError,
			expectResponse: nil,
			name:           "list images error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			tt.service.ListImages(w, r, nil)
			assert.Equal(t, tt.expectCode, w.Code)

			var actualResponse api.Response
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &actualResponse))
			if tt.expectCode == http.StatusOK {
				byteData, _ := json.Marshal(actualResponse.Data)
				var actualData []api.Img
				assert.NoError(t, json.Unmarshal(byteData, &actualData))
				assert.Equal(t, tt.expectResponse, actualData)
			} else {
				assert.NotEmpty(t, actualResponse.Errors)
			}
		})
	}
}

func TestGetImageHistory(t *testing.T) {
	imgHistory := []api.ImgHistory{
		api.ImgHistory{

			Comment:   "foo",
			CreatedBy: "bar",
			ID:        "83364c85cafc",
			RepoTags:  []string{"ref"},
			Size:      1,
		},
	}
	tests := []struct {
		service        Service
		expectCode     int
		params         httprouter.Params
		expectResponse interface{}
		name           string
	}{
		{
			service:        NewService(nil, repomocks.NewDockerRepositoryMock()),
			expectCode:     http.StatusBadRequest,
			params:         nil,
			expectResponse: imgHistory,
			name:           "image history wrong params",
		},
		{
			service: NewService(nil,
				repomocks.NewDockerRepositoryMock().
					WithImageHistory([]api.ImgHistory{}, errors.New("error"))),
			expectCode:     http.StatusInternalServerError,
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectResponse: nil,
			name:           "image history docker error",
		},
		{
			service: NewService(nil,
				repomocks.NewDockerRepositoryMock().
					WithImageHistory(imgHistory, nil)),
			expectCode:     http.StatusOK,
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectResponse: imgHistory,
			name:           "image history good",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			tt.service.GetImageHistory(w, r, tt.params)
			assert.Equal(t, tt.expectCode, w.Code)

			var actualResponse api.Response
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &actualResponse))
			if tt.expectCode == http.StatusOK {
				byteData, _ := json.Marshal(actualResponse.Data)
				var actualData []api.ImgHistory
				assert.NoError(t, json.Unmarshal(byteData, &actualData))
				assert.Equal(t, tt.expectResponse, actualData)
			} else {
				assert.NotEmpty(t, actualResponse.Errors)
			}
		})
	}
}

func TestRemoveImage(t *testing.T) {
	tests := []struct {
		service        Service
		expectCode     int
		params         httprouter.Params
		expectResponse interface{}
		name           string
	}{
		{
			service:        NewService(nil, repomocks.NewDockerRepositoryMock()),
			expectCode:     http.StatusBadRequest,
			params:         nil,
			expectResponse: nil,
			name:           "remove image wrong image id",
		},
		{
			service: NewService(nil,
				repomocks.NewDockerRepositoryMock().
					WithGetRunningContainersByImageID([]api.Ctn{api.Ctn{}}, nil)),
			expectCode:     http.StatusBadRequest,
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectResponse: nil,
			name:           "image has running containers cannot delete",
		},
		{
			service: NewService(nil,
				repomocks.NewDockerRepositoryMock().
					WithGetRunningContainersByImageID([]api.Ctn{}, nil).
					WithImageRemove(errors.New("error"))),
			expectCode:     http.StatusInternalServerError,
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectResponse: nil,
			name:           "image remove docker error",
		},
		{
			service: NewService(nil,
				repomocks.NewDockerRepositoryMock().
					WithGetRunningContainersByImageID([]api.Ctn{}, nil).
					WithImageRemove(nil)),
			expectCode:     http.StatusOK,
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectResponse: "Image was removed successfully",
			name:           "image remove ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			tt.service.RemoveImage(w, r, tt.params)
			assert.Equal(t, tt.expectCode, w.Code)

			var actualResponse api.Response
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &actualResponse))
			if tt.expectCode == http.StatusOK {
				byteData, _ := json.Marshal(actualResponse.Data)
				var actualData string
				assert.NoError(t, json.Unmarshal(byteData, &actualData))
				assert.Equal(t, tt.expectResponse, actualData)
			} else {
				assert.NotEmpty(t, actualResponse.Errors)
			}
		})
	}
}
