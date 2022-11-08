package rest

type Request struct {
	Headers  map[string]string `json:"headers,omitempty"`
	Status   int               `json:"status,omitempty"`
	Data     any               `json:"data,omitempty"`
	MetaData any               `json:"metaData,omitempty"`
}
