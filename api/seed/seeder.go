package seed

import (
	"log"

	"github.com/diptyojha/goLngFirstProject/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Nickname: "User_1",
		Email:    "User_1@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
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

var userlocations = []models.UserLocation{
	models.UserLocation{
		UserLocationID: 1,
	},
	models.UserLocation{
		UserLocationID: 2,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.UserLocation{}, &models.Location{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Location{}, &models.UserLocation{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Location{}).AddForeignKey("creator_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.UserLocation{}).AddForeignKey("creator_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		locations[i].CreatorID = users[i].ID
		userlocations[i].CreatorID = users[i].ID

		err = db.Debug().Model(&models.Location{}).Create(&locations[i]).Error
		if err != nil {
			log.Fatalf("cannot seed locations table: %v", err)
		}

		err = db.Debug().Model(&models.UserLocation{}).Create(&userlocations[i]).Error
		if err != nil {
			log.Fatalf("cannot seed userlocations table: %v", err)
		}
	}
}
