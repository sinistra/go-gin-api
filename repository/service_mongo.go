package repository

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"log"
	"os"
	"rits/carriage/api/driver"
	"rits/carriage/api/models"
)

type ServiceRepository struct{}

const ServiceDocname = "services"

// GetServices returns the list of Services
func (s ServiceRepository) GetServices() ([]models.Service, error) {
	session := driver.ConnectMongo()
	defer session.Close()

	m := session.DB(os.Getenv("MONGO_DB")).C(ServiceDocname)
	var services []models.Service
	if err := m.Find(nil).All(&services); err != nil {
		log.Println("Failed to write results:", err)
	}

	return services, nil
}

func (s ServiceRepository) GetActiveServices() ([]models.Service, error) {
	session := driver.ConnectMongo()
	defer session.Close()

	m := session.DB(os.Getenv("MONGO_DB")).C(ServiceDocname)
	var services []models.Service
	//{deleted: {$ne: true}}
	filter := bson.M{"deleted": bson.M{"$ne": true}}
	if err := m.Find(filter).All(&services); err != nil {
		log.Println("Failed to write results:", err)
	}

	return services, nil
}

func (s ServiceRepository) GetService(id string) (models.Service, error) {
	session := driver.ConnectMongo()
	defer session.Close()
	collection := session.DB(os.Getenv("MONGO_DB")).C(ServiceDocname)

	var service models.Service

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		return service, errors.New("not found")
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	if err := collection.FindId(oid).One(&service); err != nil {
		log.Println(err)
		return service, err
	}
	return service, nil
}

// AddService inserts a SERVICE in the DB
func (s ServiceRepository) AddService(service models.Service) (string, error) {
	session := driver.ConnectMongo()
	defer session.Close()
	collection := session.DB(os.Getenv("MONGO_DB")).C(ServiceDocname)

	service.ID = bson.NewObjectId()
	err := collection.Insert(service)

	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(service.ID.Hex()), nil

}

// UpdateService updates a SERVICE in the DB
func (s ServiceRepository) UpdateService(service models.Service) error {
	session := driver.ConnectMongo()
	defer session.Close()

	err := session.DB(os.Getenv("MONGO_DB")).C(ServiceDocname).UpdateId(service.ID, service)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// RemoveService deletes a SERVICE
func (s ServiceRepository) RemoveService(id string) error {
	session := driver.ConnectMongo()
	defer session.Close()

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		return errors.New("not found")
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove record
	if err := session.DB(os.Getenv("MONGO_DB")).C(ServiceDocname).RemoveId(oid); err != nil {
		return err
	}
	return nil
}
