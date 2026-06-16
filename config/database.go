package config

import (
	"fmt"
	"log"
	"os"
	"backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	var dsn string
	if dbPassword != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName,
		)
	} else {
		dsn = fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbHost, dbPort, dbName,
		)
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	err = database.AutoMigrate(
		&models.User{},
		&models.FamilyMember{},
		&models.ChatSession{},
		&models.ChatMessage{},
		&models.HospitalService{},
		&models.EmergencyRequest{},
		&models.MedicineCategory{},
		&models.Medicine{},
		&models.Prescription{},
		&models.PrescriptionItem{},
		&models.PharmacyOrder{},
		&models.PharmacyOrderItem{},
		&models.Cart{},
		&models.CartItem{},
		&models.Doctor{},
		&models.DoctorSchedule{},
		&models.Appointment{},
		&models.Polyclinic{},
		&models.PolyclinicSchedule{},
	)
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	DB = database
	fmt.Println("Database connection successfully opened")
}
