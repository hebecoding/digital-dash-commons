package rest

type Response struct {
	Status   int               `json:"status"`
	Headers  map[string]string `json:"headers,omitempty"`
	Data     any               `json:"data,omitempty"`
	Error    error             `json:"error,omitempty"`
	MetaData any               `json:"metaData,omitempty"`
}
