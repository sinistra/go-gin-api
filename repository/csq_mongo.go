package repository

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"log"
	"os"
	"rits/carriage/api/driver"
	"rits/carriage/api/models"
)

type CsqRepository struct{}

const CsqDocname = "csqs"

// GetCSQs returns the list of CSQs
func (c CsqRepository) GetCsqs() ([]models.Csq, error) {
	session := driver.ConnectMongo()
	defer session.Close()

	m := session.DB(os.Getenv("MONGO_DB")).C(CsqDocname)
	var csqs []models.Csq
	if err := m.Find(nil).All(&csqs); err != nil {
		log.Println("Failed to write results:", err)
	}

	return csqs, nil
}

func (c CsqRepository) GetActiveCsqs() ([]models.Csq, error) {
	session := driver.ConnectMongo()
	defer session.Close()

	m := session.DB(os.Getenv("MONGO_DB")).C(CsqDocname)
	var csqs []models.Csq
	//{deleted: {$ne: true}}
	filter := bson.M{"deleted": bson.M{"$ne": true}}
	if err := m.Find(filter).All(&csqs); err != nil {
		log.Println("Failed to write results:", err)
	}

	return csqs, nil
}

func (c CsqRepository) GetCsq(id string) (models.Csq, error) {
	session := driver.ConnectMongo()
	defer session.Close()
	collection := session.DB(os.Getenv("MONGO_DB")).C(CsqDocname)

	var csq models.Csq

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		return csq, errors.New("not found")
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	if err := collection.FindId(oid).One(&csq); err != nil {
		log.Println(err)
		return csq, err
	}
	return csq, nil
}

// AddCsq inserts a CSQ in the DB
func (c CsqRepository) AddCsq(csq models.Csq) (string, error) {
	session := driver.ConnectMongo()
	defer session.Close()
	collection := session.DB(os.Getenv("MONGO_DB")).C(CsqDocname)

	csq.ID = bson.NewObjectId()
	err := collection.Insert(csq)

	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(csq.ID.Hex()), nil

}

// UpdateCsq updates a CSQ in the DB
func (c CsqRepository) UpdateCsq(csq models.Csq) error {
	session := driver.ConnectMongo()
	defer session.Close()

	err := session.DB(os.Getenv("MONGO_DB")).C(CsqDocname).UpdateId(csq.ID, csq)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// RemoveCsq deletes a CSQ
func (c CsqRepository) RemoveCsq(id string) error {
	session := driver.ConnectMongo()
	defer session.Close()

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		return errors.New("not found")
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove record
	if err := session.DB(os.Getenv("MONGO_DB")).C(CsqDocname).RemoveId(oid); err != nil {
		return err
	}
	return nil
}
