package rest

type Response struct {
	Headers  map[string]string `json:"headers,omitempty"`
	Status   int               `json:"status"`
	Error    string            `json:"error,omitempty"`
	Data     any               `json:"data,omitempty"`
	MetaData any               `json:"metaData,omitempty"`
}
