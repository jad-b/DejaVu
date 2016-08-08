package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
)

type FileMemory struct {
	Scope     string
	Directory string
	filename  string
	cache     *http.Response
}

func NewFileMemory(scope, directory string) *FileMemory {
	return &FileMemory{
		Directory: directory,
		Scope:     scope,
		filename:  path.Join(fm.Directory, fm.scope+".json"),
	}
}

// Recall returns a file based on its scope.
func (fm *FileMemory) Recall(_ *http.Request) *http.Response {
	if cache != nil {
		return fm.cache
	}
	f, err := os.Open(fm.filename)
	if err != nil {
		log.Printf("No recorded response(s) found at %s", filename)
		return nil
	}
	resp, err := http.ReadResponse(bufio.NewReader(f), nil)
	if err != nil {
		log.Printf("Failed to read response for scope '%s'", fm.scope)
		return nil
	}
	fm.cache = resp
	return resp
}

func (fm *FileMemory) Remember(r *http.Request, resp *http.Response) error {
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	f, err := os.Create(fm.filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, bufio.NewReader(b)); err != nil {
		return err
	}
	return nil
}
