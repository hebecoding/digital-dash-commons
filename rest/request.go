package rest

type Request struct {
	Headers  map[string]string `json:"headers,omitempty"`
	Data     any               `json:"data"`
	MetaData any               `json:"metaData,omitempty"`
}
