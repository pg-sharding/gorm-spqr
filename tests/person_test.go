package tests

import (
	"fmt"
	"github.com/pg-sharding/gorm-spqr/controllers"
	"github.com/pg-sharding/gorm-spqr/models"
	"github.com/pg-sharding/gorm-spqr/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func setup() {
	models.ConnectDatabase()
}

func tearDown(t *testing.T) {
	if err := utils.DeleteAll(); err != nil {
		t.Errorf("Error running tearDown: %s", err)
	}
}

func TestMain(m *testing.M) {
	models.SetupSharding()
	code := m.Run()
	_ = utils.DeleteAll()
	os.Exit(code)
}

func TestCreatePerson(t *testing.T) {
	setup()

	var person = models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
	}

	if err := controllers.WritePerson(&person); err != nil {
		t.Errorf("could not write: %s", err)
	}

	var allPersons = controllers.GetAllPersons()
	if len(allPersons) != 1 {
		t.Errorf("Expected to have 1 person in db, got %d", len(allPersons))
	}

	tearDown(t)
}

func TestReadPerson(t *testing.T) {
	setup()

	var person = models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
	}

	if err := controllers.WritePerson(&person); err != nil {
		t.Errorf("could not write: %s", err)
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

func TestUpdatePerson(t *testing.T) {
	setup()

	var person = models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
	}

	if err := controllers.WritePerson(&person); err != nil {
		t.Errorf("could not write: %s", err)
	}

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
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
	}

	if err := controllers.WritePerson(&person); err != nil {
		t.Errorf("could not write: %s", err)
	}

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

func TestMultipleWrite(t *testing.T) {
	setup()

	var firstShard = make([]*models.Person, 0)
	for i := 1; i < 100; i++ {
		person := &models.Person{
			ID:        uint32(i),
			FirstName: fmt.Sprintf("Name_%d", i),
		}
		firstShard = append(firstShard, person)
	}

	if err := controllers.WritePeople(firstShard); err != nil {
		t.Errorf("could not write: %s", err)
	}

	var secondShard = make([]*models.Person, 0)
	for i := 100; i < 200; i++ {
		person := &models.Person{
			ID:        uint32(i),
			FirstName: fmt.Sprintf("Name_%d", i),
		}
		secondShard = append(secondShard, person)
	}

	if err := controllers.WritePeople(secondShard); err != nil {
		t.Errorf("could not write: %s", err)
	}

	firstShardPeople, err := controllers.GetPeople(uint32(1), uint32(99))
	if err != nil {
		t.Errorf("could not get people: %s", err)
	}

	assert.Equalf(t, len(firstShardPeople), 99, "expected to have %d records on 1st shard, got %d", 99, len(firstShardPeople))

	secondShardPeople, err := controllers.GetPeople(uint32(100), uint32(199))
	if err != nil {
		t.Errorf("could not get people: %s", err)
	}

	assert.Equalf(t, len(secondShardPeople), 100, "expected to have %d records on 1st shard, got %d", 100, len(secondShardPeople))
}
