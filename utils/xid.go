package utils

import (
	"encoding/json"

	"github.com/rs/xid"
)

type XID struct {
	ID string `json:"_id" bson:"_id"`
}

func NewXID() XID {
	return XID{xid.New().String()}
}
func (id *XID) String() string {
	return id.ID
}

func (id XID) MarshalBSON() ([]byte, error) {
	return json.Marshal(id.ID)
}

func (id *XID) UnmarshalBSON(data []byte) error {
	var xid string

	if err := json.Unmarshal(data, &xid); err != nil {
		return err
	}

	id.ID = xid

	return nil
}
