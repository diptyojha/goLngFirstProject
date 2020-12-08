package controllertests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/diptyojha/goLngFirstProject/api/controllers"
	"github.com/diptyojha/goLngFirstProject/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var userInstance = models.User{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		models.User{
			Nickname: "Grand",
			Email:    "grand@gmail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}
	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserAndLocationTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Location{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Location{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneLocation() (models.Location, error) {

	err := refreshUserAndLocationTable()
	if err != nil {
		return models.Location{}, err
	}
	user := models.User{
		Nickname: "User_111",
		Email:    "user_111@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Location{}, err
	}
	location := models.Location{
		Loc_Name: "Location 1",
		Address:  "Hello world 1",
		Pincode:  "12345",
		TimeZone: "EST",
	}

	err = server.DB.Model(&models.Location{}).Create(&location).Error
	if err != nil {
		return models.Location{}, err
	}
	return location, nil
}

func seedUsersAndLocations() ([]models.User, []models.Location, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Location{}, err
	}
	var users = []models.User{
		models.User{
			Nickname: "User_1",
			Email:    "User_1@gmail.com",
			Password: "password",
		},
		models.User{
			Nickname: "User_2",
			Email:    "User_2@gmail.com",
			Password: "password",
		},
	}
	var locations = []models.Location{
		models.Location{
			Loc_Name: "Location 1",
			Address:  "Hello world 1",
			Pincode:  "12345",
			TimeZone: "EST",
		},
		models.Location{
			Loc_Name: "Location 2",
			Address:  "Hello world 2",
			Pincode:  "126566",
			TimeZone: "EST",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		locations[i].CreatorID = users[i].ID

		err = server.DB.Model(&models.Location{}).Create(&locations[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, locations, nil
}
