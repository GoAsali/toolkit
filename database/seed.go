package database

import (
	"gorm.io/gorm"
)

type ISeedDatabase interface {
	Seed(db *gorm.DB) error
}

func Seed(seeders ...ISeedDatabase) error {
	database, err := Database()
	if err != nil {
		return err
	}
	for _, seed := range seeders {
		if err := seed.Seed(database); err != nil {
			return err
		}
	}

	return nil
}
