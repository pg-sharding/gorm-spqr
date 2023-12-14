package tests

import (
	"github.com/pg-sharding/gorm-spqr/controllers"
	"github.com/pg-sharding/gorm-spqr/models"
	"github.com/pg-sharding/gorm-spqr/utils"
	"reflect"
	"testing"
	"time"
)

func setup() {
	models.ConnectDatabase()
}

func tearDown(t *testing.T) {
	if err := utils.DeleteAll(); err != nil {
		t.Errorf("Error running tearDown: %s", err)
	}
}

func TestCreatePerson(t *testing.T) {
	setup()

	var person = models.Person{
		ID:         1,
		FirstName:  "John",
		LastName:   "Smith",
		JoinedDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0)),
	}

	controllers.WritePerson(&person)

	var allPersons = controllers.GetAllPersons()
	if len(allPersons) != 1 {
		t.Errorf("Expected to have 1 person in db, got %d", len(allPersons))
	}

	tearDown(t)
}

func TestReadPerson(t *testing.T) {
	setup()

	var person = models.Person{
		ID:         1,
		FirstName:  "John",
		LastName:   "Smith",
		JoinedDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0)),
	}

	controllers.WritePerson(&person)

	var personDb *models.Person
	var err error
	if personDb, err = controllers.GetPerson(person.ID); err != nil {
		t.Errorf("error getting person: %s", err)
	}
	if !reflect.DeepEqual(person, *personDb) {
		t.Errorf("Expected %#v, got %#v", person, personDb)
	}

	tearDown(t)
}

func TestUpdatePerson(t *testing.T) {
	setup()

	var person = models.Person{
		ID:         1,
		FirstName:  "John",
		LastName:   "Smith",
		JoinedDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0)),
	}

	controllers.WritePerson(&person)

	person.Email = "foo@bar"

	if err := controllers.UpdatePerson(&person); err != nil {
		t.Errorf("error updating: %s", err)
	}

	var personDb *models.Person
	var err error
	if personDb, err = controllers.GetPerson(person.ID); err != nil {
		t.Errorf("error getting person: %s", err)
	}
	if !reflect.DeepEqual(person, *personDb) {
		t.Errorf("Expected %#v, got %#v", person, personDb)
	}

	tearDown(t)
}

func TestDeletePerson(t *testing.T) {
	setup()

	var person = models.Person{
		ID:         1,
		FirstName:  "John",
		LastName:   "Smith",
		JoinedDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0)),
	}

	controllers.WritePerson(&person)

	var allPersons = controllers.GetAllPersons()
	if len(allPersons) != 1 {
		t.Errorf("Expected to have 1 person in db, got %d", len(allPersons))
	}

	if err := controllers.DeletePerson(person.ID); err != nil {
		t.Errorf("error deleting person: %s", err)
	}

	allPersons = controllers.GetAllPersons()
	if len(allPersons) != 0 {
		t.Errorf("Expected to have no people in db, got %d", len(allPersons))
	}

	tearDown(t)
}
