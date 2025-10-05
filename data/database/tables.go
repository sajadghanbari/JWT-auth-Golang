package db

import "JWT-Authentication-go/data/models"

func Tables() []interface{} {
	return []interface{}{
		&models.User{},
		// &models.Product{},
		// &models.Order{},
	}
}