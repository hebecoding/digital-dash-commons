package utils

import (
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
)

type XID string

func NewXID() XID {
	return XID(xid.New().Bytes())
}
func (id *XID) String() string {
	return string(*id)
}

func (id XID) MarshalBSON() ([]byte, error) {
	return bson.Marshal(struct {
		ID string
	}{string(id)})
}

func (id *XID) UnmarshalBSON(data []byte) error {
	var tmp struct{ ID string }

	if err := bson.Unmarshal(data, &tmp); err != nil {
		return err
	}

	*id = XID(tmp.ID)

	return nil
}
