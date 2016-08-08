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

type TreeMemory struct {
	filename string
	Tree     interface{}
}

func NewTreeMemory(filename string) *TreeMemory {
	return &TreeMemory{
		filename: filename,
	}
}

// Recall returns the Response matching the Request.
func (tm *TreeMemory) Recall(r *http.Request) {
	return nil
}

func (tm *TreeMemory) Remember(r *http.Request, resp *http.Response) error {
	// Dump to bytes, preserving Body for reading by caller
	b, dumpErr := httputil.DumpResponse(resp, true)
	if dumpErr != nil {
		log.Print("Failed to dump response")
		return err
	}
	// Read bytes back into http.Response
	respCopy, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(b)), nil)
	if err != nil {
		log.Print(err)
		return err
	}
	//tm.Tree = respCopy // Save our copy into the given map
	return nil
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
