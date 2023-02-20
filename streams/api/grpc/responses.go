package grpc

import "time"

type oneRes struct {
	id          string
	name        string
	owner       string
	url         string
	price       uint64
	external    bool
	offer       bool
	project     string
	dataset     string
	table       string
	fields      string
	terms       string
	visibility  string
	accessType  string
	maxCalls    uint64
	maxUnit     string
	endDate     *time.Time
	err         error
	subCategory string
}
