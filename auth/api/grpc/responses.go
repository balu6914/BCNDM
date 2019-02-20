package grpc

type identityRes struct {
	id  string
	err error
}

type emailRes struct {
	email        string
	contactEmail string
	err          error
}

type partnersRes struct {
	partners []string
	err      error
}
