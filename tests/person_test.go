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

	err := controllers.WritePerson(&person)
	assert.NoErrorf(t, err, "could not write: %s", err)

	var allPersons = controllers.GetAllPersons()
	assert.Len(t, allPersons, 1, "Expected to have 1 person in db, got %d", len(allPersons))

	tearDown(t)
}

func TestReadPerson(t *testing.T) {
	setup()

	var person = models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
	}

	err := controllers.WritePerson(&person)
	assert.NoErrorf(t, err, "could not write: %s", err)

	var personDb *models.Person
	personDb, err = controllers.GetPerson(person.ID)
	assert.NoErrorf(t, err, "error getting person: %s", err)
	assert.True(t, reflect.DeepEqual(person, *personDb), "Expected %#v, got %#v", person, personDb)

	tearDown(t)
}

func TestUpdatePerson(t *testing.T) {
	setup()

	var person = models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
	}

	err := controllers.WritePerson(&person)
	assert.NoErrorf(t, err, "could not write: %s", err)

	person.Email = "foo@bar"

	err = controllers.UpdatePerson(&person)
	assert.NoErrorf(t, err, "error updating: %s", err)

	var personDb *models.Person
	personDb, err = controllers.GetPerson(person.ID)
	assert.NoErrorf(t, err, "error getting person: %s", err)
	assert.True(t, reflect.DeepEqual(person, *personDb), "Expected %#v, got %#v", person, personDb)

	tearDown(t)
}

func TestDeletePerson(t *testing.T) {
	setup()

	var person = models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
	}

	err := controllers.WritePerson(&person)
	assert.NoErrorf(t, err, "could not write: %s", err)

	var allPersons = controllers.GetAllPersons()
	assert.Len(t, allPersons, 1, "Expected to have 1 person in db, got %d", len(allPersons))

	err = controllers.DeletePerson(person.ID)
	assert.NoErrorf(t, err, "error deleting person: %s", err)

	allPersons = controllers.GetAllPersons()
	assert.Emptyf(t, allPersons, "Expected to have no people in db, got %d", len(allPersons))

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

	err := controllers.WritePeople(firstShard)
	assert.NoErrorf(t, err, "could not write: %s", err)

	var secondShard = make([]*models.Person, 0)
	for i := 100; i < 200; i++ {
		person := &models.Person{
			ID:        uint32(i),
			FirstName: fmt.Sprintf("Name_%d", i),
		}
		secondShard = append(secondShard, person)
	}

	err = controllers.WritePeople(secondShard)
	assert.NoErrorf(t, err, "could not write: %s", err)

	firstShardPeople, err := controllers.GetPeople(uint32(1), uint32(99))
	assert.NoErrorf(t, err, "could not get people: %s", err)

	assert.Equalf(t, len(firstShardPeople), 99, "expected to have %d records on 1st shard, got %d", 99, len(firstShardPeople))

	secondShardPeople, err := controllers.GetPeople(uint32(100), uint32(199))
	assert.NoErrorf(t, err, "could not get people: %s", err)

	assert.Equalf(t, len(secondShardPeople), 100, "expected to have %d records on 1st shard, got %d", 100, len(secondShardPeople))
}
