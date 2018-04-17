package api

type saveStreamReq struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Price       int    `json:"price"`
	// owner       User
	// longlat     Location
}

type oneStreamReq struct {
	Name string
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

type removeStreamReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type purchaseStreamReq struct {
	Name  string `json:"name"`
	Hours int    `json:"hours"`
}
