package api

type modifyStreamReq struct {
	Id          string
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	Price       int    `json:"price,omitempty"`
	// owner       User
	// longlat     Location
}

type readStreamReq struct {
	Id string
}

type searchStreamReq struct {
	Type string
	x0   int
	y0   int
	x1   int
	y1   int
	x2   int
	y2   int
	x3   int
	y3   int
}

type purchaseStreamReq struct {
	Id    string
	Hours int `json:"hours"`
}
