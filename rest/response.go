package rest

type Response struct {
	Status   int `json:"status"`
	Data     any `json:"data,omitempty"`
	MetaData any `json:"metaData,omitempty"`
}

type ErrResponse struct {
	Status   int    `json:"status"`
	Data     any    `json:"data,omitempty"`
	Error    string `json:"error,omitempty"`
	MetaData any    `json:"metaData,omitempty"`
}
