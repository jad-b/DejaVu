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

// Memory stores and retrieves remembered Request/Response pairs.
type Memory interface {
	Recall(*http.Request) *http.Response
	Remember(*http.Request, *http.Response)
	Persist() error
}

type TreeMemory struct {
	filename string
	Tree     struct{}
}

func NewTreeMemory(filename string) *TreeMemory {
	return &TreeMemory{
		filename: filename,
		Tree: struct{}
	}
}

// Recall returns the Response matching the Request.
func (tm *TreeMemory) Recall(r *http.Request) {
	return nil
}

func (tm *TreeMemory) Remember(r *http.Request, resp *http.Response) {
	// Dump to bytes, preserving Body for reading by caller
	b, dumpErr := httputil.DumpResponse(resp, true)
	if dumpErr != nil {
		log.Print("Failed to dump response")
		return resp, err
	}
	// Read bytes back into http.Response
	respCopy, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(b)), nil)
	if err != nil {
		log.Print(err)
		return resp, err
	}
	//tm.Tree = respCopy // Save our copy into the given map
}

// Persist stores the Tree to a file.
func (tm *TreeMemory) Persists() error {
	f, err := os.Create(tm.filename)
	if err != nil {
		return err
	}
	// TODO Save responses w/ Body to JSON map.
	b, err := json.MarshalIndent(tm.Tree, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}
