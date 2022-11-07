package main

import (
	"database/sql"
	"fmt"
	"time"
)

type DBManager struct {
	db *sql.DB
}

func NewDBManager(db *sql.DB) DBManager {
	return DBManager{db}
}

type Car struct {
	ID           int64
	Manufacturer string
	Model        string
	Year         int64
	Typeof_car   string
	Price        float64
	Mileage      int64
	Image_url    string
	Created_at   time.Time
	Images       []*CarImage
}

type CarImage struct {
	ID             int64
	ImageUrl       string
	SequenceNumber int32
}

type GetCarsParams struct {
	Limit  int32
	Page   int32
	Search string
}

type GetCarResponse struct {
	Cars  []*Car
	Count int32
}

func (c *DBManager) CreateCar(car *Car) (int64, error) {
	var carID int64

	query := `
		INSERT INTO cars (
			id,
			manufacturer,
			model,
			year,
			type_car,
			price,
			mileage,
			image_url
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	row := c.db.QueryRow(
		query,
		car.ID,
		car.Manufacturer,
		car.Model,
		car.Year,
		car.Typeof_car,
		car.Price,
		car.Mileage,
		car.Image_url,
	)
	err := row.Scan(&carID)
	if err != nil {
		return 0, err
	}

	queryInsertImage := `
		INSERT INTO cars_images (
			car_id,
			image_url,
			sequence_number
		) VALUES($1, $2, $3)
	`

	for _, image := range car.Images {
		_, err := c.db.Exec(
			queryInsertImage,
			carID,
			image.ImageUrl,
			image.SequenceNumber,
		)
		if err != nil {
			return 0, err
		}
	}

	return carID, err

}

func (c *DBManager) GetCar(id int64) (*Car, error) {
	var car Car

	car.Images = make([]*CarImage, 0)

	query := `
		SELECT
			c.id,
			c.manufacturer
			c.model
			c.year
			c.type_car
			c.price
			c.mileage
			c.image_url
		FROM cars c
		WHERE c.id=$1
	`

	row := c.db.QueryRow(query, id)
	err := row.Scan(
		&car.ID,
		&car.Manufacturer,
		&car.Model,
		&car.Year,
		&car.Typeof_car,
		&car.Price,
		&car.Mileage,
		&car.Image_url,
		&car.Created_at,
	)
	if err != nil {
		return nil, err
	}

	queryImages := `
		SELECT 
			id,
			images_url,
			sequence_number
		FROM cars_images
		WHERE car_id=$1
	`

	rows, err := c.db.Query(queryImages, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var image CarImage

		err := rows.Scan(
			&image.ID,
			&image.ImageUrl,
			&image.SequenceNumber,
		)
		if err != nil {
			return nil, err
		}
		car.Images = append(car.Images, &image)
	}
	fmt.Println(&car)
	return &car, nil

}

func (c *DBManager) GetAllCars(params *GetCarsParams) (*GetCarResponse, error) {
	
	var result GetCarResponse

	result.Cars = make([]*Car, 0)

	filter := ""
	if params.Search != "" {
		filter = fmt.Sprintf("WHERE name ilike '%s'", "%"+params.Search+"%")
	}

	query := `
		SELECT
		c.id,
		c.manufacturer
		c.model
		c.year
		c.type_car
		c.price
		c.mileage
		c.image_url
		c.created_at
		FROM cars c
		WHERE c.id=$1
	` + filter + `
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`
	offset := (params.Page - 1) * params.Limit
	rows, err := c.db.Query(query, params.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var car Car

		err := rows.Scan(
			&car.ID,
			&car.Manufacturer,
			&car.Model,
			&car.Year,
			&car.Typeof_car,
			&car.Price,
			&car.Mileage,
			&car.Image_url,
			&car.Created_at,
		)
		if err != nil {
			return nil, err
		}
		result.Cars = append(result.Cars, &car)
	}
	return &result, nil

}



func (c *DBManager) UpdateProduct(car *Car) error {
	query := `
		UPDATE cars SET
			manufacturer=$1
			model=$2
			year=$3
			type_car=$4
			price=$5
			mileage=$6
			image_url=$7
		WHERE id=$8		
	`

	result, err := c.db.Exec(
		query,
		car.Manufacturer,
		car.Model,
		car.Year,
		car.Typeof_car,
		car.Price,
		car.Mileage,
		car.Image_url,
	)

	if err != nil {
		return err
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		return sql.ErrNoRows
	}
	if rowsCount == 0 {
		return sql.ErrNoRows
	}

	queryDeleteImages := `DELETE FROM  cars_images WHERE
	car_id=$1`
	_, err = c.db.Exec(queryDeleteImages, car.ID)
	if err != nil {
		return err
	}

	queryInsertImage := `
		INSERT INTO cars_images (
			car_id,
			images_url,
			sequence_number
		) VALUES ($1, $2, $3)
	`

	for _, image := range car.Images {
		_, err := c.db.Exec(
			queryInsertImage,
			image.SequenceNumber,
		)
		if err != nil {
			return nil
		}
	}

	return nil
}