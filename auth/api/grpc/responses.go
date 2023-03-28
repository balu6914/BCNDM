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

type existsRes struct {
	err error
}

type userResp struct {
	id           string
	email        string
	contactEmail string
	firstName    string
	lastName     string
	company      string
	address      string
	country      string
	mobile       string
	phone        string
	role         string
	err          error
}
