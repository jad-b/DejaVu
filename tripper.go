package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (rt RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

// DejaVu saves HTTP responses by Request.
//
// When Record is false, DejaVu replays responses by matching against requests.
type DejaVu struct {
	Record      bool
	Tripper     http.RoundTripper
	StorageLocation *url.URL
	memory Memory
}

// NewDejaVu sets up the DejaVu object. It initializes to neither recording nor
// replaying, so this needs to be set explicitly after creation.
func NewDejaVu(filename string, rt http.RoundTripper) *DejaVu {
	if rt == nil {
		rt = http.DefaultTransport
	}
	earl, err := url.Parse(filename)
	if err != nil {
		panic(err)
	}
	if earl.Scheme == "" {
		earl.Scheme = "file"
	}
	return &DejaVu{
		Record: true,
		StorageLocation:    earl,
		Tripper:     rt,
		memory: NewTreeMemory(earl.Path),
	}
}

// Client returns a *http.Client with the the Transport RoundTripper set to the
// configurd DejaVu RoundTripper.
func (dv *DejaVu) NewClient() *http.Client {
	return &http.Client{
		Transport: dv.Tripper,
	}
}

// WrapClient repalces the *http.Client's Transport with a closure that
// records or replays the response.
func (dv *DejaVu) WrapClient(c *http.Client) *http.Client {
	// Wrap the old Tranport inside DejaVu
	dv.Tripper := client.Transport
	// Insert our wrapper into the Client
	client.Transport = dv
	return client
}

// RoundTrip executes the underlying RoundTripper, then saves the response.
func (dv *DejaVu) RoundTrip(r *http.Request) (*http.Response, error) {
	if !dv.Record {
		return dv.Storage.Recall(r), nil
	}

	// Execute request, get response
	resp, err := rt.Tripper.RoundTrip(r)
	if err != nil {
		return resp, err
	}

	if dv.Record {
		dv.Storage.Remember(r, resp)
	}

	// Forward original response
	return resp, err
}
