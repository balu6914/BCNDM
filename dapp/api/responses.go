package api

import (
	"monetasa/dapp"
)

type statusRes struct {
	Status string `json:"status"`
}

type saveStreamRes struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Price       int    `json:"price"`
	// owner       User
	// longlat     Location
}

type oneStreamRes struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type searchStreamRes struct {
	// Streams []Stream `json:"streams"`
	Streams []Stream
}

type removeStreamRes struct {
	Status string `json:"status"`
}
