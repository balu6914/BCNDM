package api

import (
	"monetasa/dapp"
)

const contentType = "application/json; charset=utf-8"

type versionRes struct {
	Version string `json:"version"`
}

type modifyStreamRes struct {
	Status string `json:"status"`
}

type readStreamRes struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type searchStreamRes struct {
	// Streams []Stream `json:"streams"`
	Streams []dapp.Stream
}
