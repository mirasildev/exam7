package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	Host     = "localhost"
	User     = "postgres"
	Port     = 5432
	Password = "1105"
	Database = "exam_project"
)

func main() {
	connStr := fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s
	sslmode=disable`,
		Host,
		Port,
		User,
		Password,
		Database,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open connection: %v", err)
	}

	c := NewDBManager(db)

	id, err := c.CreateCar(&Car{
		ID: 2,
		Manufacturer: "BMW",
		Model:        "X5",
		Year:         2020,
		Typeof_car:   "SUV",
		Price:        50000,
		Mileage:      2000,
		Image_url:    "test_url",
		Images: []*CarImage{
			{
				ImageUrl:       "test_url1",
				SequenceNumber: 1,
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to create car: %v", err)
	}

	car1, err := c.GetCar(id)
	if err != nil {
		log.Fatalf("failed to get a car: %v", err)
	}
	fmt.Println(car1)

	resp, err := c.GetAllCars(&GetCarsParams{
		Limit: 10,
		Page:  1,
	})

	if err != nil {
		log.Fatalf("failed to get all cars: %v", err)
	}

	fmt.Printf("%v", resp)

	
}
