package main

import (
	"advertiser/shared/pkg/service/repo"
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
)

func main() {
	db := repo.New()

	topics := []models.Topic{
		{ID: "art"},
		{ID: "books"},
		{ID: "food"},
		{ID: "pets"},
		{ID: "sport"},
	}

	for _, topic := range topics {
		result := db.Create(&topic)
		if result.Error != nil {
			fmt.Printf("failed to create topic %s: %s\n", topic.ID, result.Error)
		}
	}
}
