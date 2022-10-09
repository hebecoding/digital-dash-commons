package rest

type Response struct {
	Status   int `json:"status"`
	MetaData any `json:"metaData,omitempty"`
	Data     any `json:"data,omitempty"`
}
