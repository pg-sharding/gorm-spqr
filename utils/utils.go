package utils

import "github.com/pg-sharding/gorm-spqr/models"

func DeleteAll() error {
	models.DB.Where("1 <= id AND id <= 99").Delete(&models.Person{})
	models.DB.Where("100 <= id AND id <= 199").Delete(&models.Person{})
	return nil
}
