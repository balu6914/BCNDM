package grpc

type oneRes struct {
	id       string
	name     string
	owner    string
	url      string
	price    uint64
	external bool
	project  string
	dataset  string
	table    string
	fields   string
	terms    string
	err      error
}
