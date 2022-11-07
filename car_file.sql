CREATE TABLE "cars"(
    "id" INTEGER PRIMARY KEY,
    "manufacturer" VARCHAR(255),
    "model" VARCHAR(255),
    "year" INTEGER,
    "type_car" VARCHAR(255),
    "price" DECIMAL(18,2),
    "mileage" INTEGER,
    "image_url" VARCHAR,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
  ------------- ------------
CREATE TABLE "cars_images"(
    "id" INTEGER PRIMARY KEY,
    "images_url" VARCHAR,
    "sequence_number" INTEGER,
    "car_id" INTEGER,
    CONSTRAINT cars_id_foreign
        FOREIGN KEY(id)
            REFERENCES cars(id)
            on delete cascade
);

