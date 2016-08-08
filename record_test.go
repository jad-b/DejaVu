package main

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	tasks "google.golang.org/api/tasks/v1"
)

var (
	apiClient *http.Client
	apiURL    = &url.URL{
		Scheme: "https",
		Host:   "www.googleapis.com",
	}
)

func init() {
	// Setup the Goole OAuth API Client
	apiClient := loadClient("./client_secret.json", tasks.TasksReadonlyScope)
}

// Must panics if an error is encountered during Request creation.
func Must(r *http.Request, err error) *http.Request {
	if err != nil {
		panic(err)
	}
	return r
}

func TestRecording(t *testing.T) {
	dv := NewDejaVu()
	dv.Record = true
	client = dv.WrapClient(apiClient)
	earl := *apiURL
	earl.Path = tc.Path
	req := Must(http.NewRequest(http.MethodGet, earl.String(), nil))
	resp, err := client.Do(dv)
	if err != nil {
		t.Error(err)
		continue
	}
	// Now, test replaying this transaction
	dv.Record = false
	recResp, err = client.Do(dv)
	if !reflect.DeepEqual(resp, recResp) {
		t.Fatalf("Recorded response did not equal returned response\n: %v != %v",
			recResp, resp)
	}
}

func TestResponseRecording(t *testing.T) {
	dv := NewDejaVu()
	dv.Record = true
	client = dv.WrapClient(apiClient)
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
		earl := *apiURL
		earl.Path = tc.Path
		req := Must(http.NewRequest(tc.Method, earl.String(), nil))
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
			continue
		}
		// Verify the Response got saved
		if resp = dv.Recall(req); resp == nil {
			t.Errorf("No saved Response to match the Request")
			continue
		}
	}
}
