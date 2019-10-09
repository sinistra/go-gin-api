package models

import (
	"github.com/globalsign/mgo/bson"
)

type Service struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string        `json:"name,omitempty" bson:"name,omitempty"`
	Vendor    string        `json:"vendor,omitempty" bson:"vendor,omitempty"`
	Product   string        `json:"product,omitempty" bson:"product,omitempty"`
	Bandwidth string        `json:"bandwidth,omitempty" bson:"bandwidth,omitempty"`
	Deleted   bool          `json:"deleted,omitempty" bson:"deleted,omitempty"`
}
