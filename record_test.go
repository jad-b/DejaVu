package main

import (
	"net/http"
	"net/url"
	"testing"
)

var apiURL = &url.URL{
	Scheme: "https",
	Host:   "www.googleapis.com",
}

// Must panics if an error is encountered during Request creation.
func Must(r *http.Request, err error) *http.Request {
	if err != nil {
		panic(err)
	}
	return r
}

func TestResponseRecording(t *testing.T) {
	dv := NewDejaVu()
	dv.Record = true
	client := loadClient() // Complete OAuth login to Google
	client = dv.WrapClient(client)
	testCases := []struct {
		Method   string
		Path     string
		RespCode int
	}{
		{
			Method:   http.MethodGet,
			Path:     "/tasks/v1/users/@me/lists",
			RespCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		// Execute a Request
		earl := *apiURL     // Copy the base URL
		earl.Path = tc.Path // Insert the Path
		resp, err := client.Do(Must(http.NewRequest(tc.Method, earl.String(), nil)))
		if err != nil {
			t.Error(err)
			continue
		}
		// Verify the Response got saved
		if _, ok := dv.ResponseMap[earl.Path]; !ok {
			t.Errorf("Missing path '%s' from the saved responses", earl.Path)
			continue
		}
		if _, ok := dv.ResponseMap[earl.Path][tc.Method]; !ok {
			t.Errorf("Missing method '%s' from ResponseMap[%s]", tc.Method, earl.Path)
			continue
		}
	}
}
