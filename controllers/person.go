package controllers

import "github.com/pg-sharding/gorm-spqr/models"

func GetAllPersons() []models.Person {
	var res []models.Person
	models.DB.Find(&res)

	return res
}

func GetPerson(id uint32) (*models.Person, error) {
	var person models.Person
	if err := models.DB.Where("id = ?", id).First(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func WritePerson(person *models.Person) {
	models.DB.Create(person)
}

func UpdatePerson(person *models.Person) error {
	var current models.Person
	if err := models.DB.Where("id = ?", person.ID).First(&current).Error; err != nil {
		return err
	}
	models.DB.Save(&person)
	return nil
}

func DeletePerson(id uint32) error {
	var person models.Person
	if err := models.DB.Where("id = ?", id).First(&person).Error; err != nil {
		return err
	}

	models.DB.Delete(&person)
	return nil
}
