package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Csq struct {
	ID           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerName string        `json:"customer_name,omitempty" bson:"customer_name,omitempty"`
	Service      string        `json:"service,omitempty" bson:"service,omitempty"`
	Bandwidth    string        `json:"bandwidth,omitempty" bson:"bandwidth,omitempty"`
	Locations    []Location    `json:"locations,omitempty" bson:"locations,omitempty"`
	Status       string        `json:"status,omitempty" bson:"status,omitempty"`
	CreatedBy    string        `json:"created_by,omitempty" bson:"created_by,omitempty"`
	RequestDate  time.Time     `json:"request_date" bson:"request_date"`
	UpdatedBy    string        `json:"updated_by" bson:"updated_by"`
	LastUpdate   time.Time     `json:"lastupdate" bson:"lastupdate"`
	Deleted      bool          `json:"deleted,omitempty" bson:"deleted,omitempty"`
}

type Location struct {
	Name             string  `json:"name,omitempty" bson:"name,omitempty"`
	FormattedAddress string  `json:"formatted_address,omitempty" bson:"formatted_address,omitempty"`
	Latitude         float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude        float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
	PlaceID          string  `json:"place_id,omitempty" bson:"place_id,omitempty"`
	Status           string  `json:"status,omitempty" bson:"status,omitempty"`
	Product          string  `json:"product,omitempty" bson:"product,omitempty"`
	Service          string  `json:"service,omitempty" bson:"service,omitempty"`
}
