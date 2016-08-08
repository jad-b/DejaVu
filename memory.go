package main

import "net/http"

// Memory stores and retrieves remembered Request/Response pairs.
type Memory interface {
	Recall(*http.Request) *http.Response
	Remember(*http.Request, *http.Response) error
	Persist() error
}
