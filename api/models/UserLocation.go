package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type UserLocation struct {
	ID             uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserLocation   Location  `json:"user_location"`
	UserLocationID uint64    `gorm:"not null" json:"user_location_id"`
	Creator        User      `json:"creator"`
	CreatorID      uint32    `gorm:"not null" json:"creator_id"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *UserLocation) Prepare() {
	p.ID = 0
	p.UserLocation = Location{}
	p.Creator = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *UserLocation) Validate() error {

	if p.UserLocationID < 1 {
		return errors.New("Required UserLocationID")
	}
	if p.CreatorID < 1 {
		return errors.New("Required CreatorID")
	}
	return nil
}

func (p *UserLocation) SaveUserLocation(db *gorm.DB) (*UserLocation, error) {
	var err error
	err = db.Debug().Model(&UserLocation{}).Create(&p).Error
	if err != nil {
		return &UserLocation{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorID).Take(&p.Creator).Error
		if err != nil {
			return &UserLocation{}, err
		}
	}
	if p.UserLocationID != 0 {
		err = db.Debug().Model(&Location{}).Where("id = ?", p.UserLocationID).Take(&p.UserLocation).Error
		if err != nil {
			return &UserLocation{}, err
		}
	}
	return p, nil
}

func (p *UserLocation) FindAllUserLocations(db *gorm.DB) (*[]UserLocation, error) {
	var err error
	userLocations := []UserLocation{}
	err = db.Debug().Model(&UserLocation{}).Limit(100).Find(&userLocations).Error
	if err != nil {
		return &[]UserLocation{}, err
	}
	if len(userLocations) > 0 {
		for i, _ := range userLocations {
			err := db.Debug().Model(&User{}).Where("id = ?", userLocations[i].CreatorID).Take(&userLocations[i].Creator).Error
			if err != nil {
				return &[]UserLocation{}, err
			}
		}
	}
	return &userLocations, nil
}

func (p *UserLocation) FindUserLocationByID(db *gorm.DB, pid uint64) (*UserLocation, error) {
	var err error
	err = db.Debug().Model(&UserLocation{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &UserLocation{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorID).Take(&p.Creator).Error
		if err != nil {
			return &UserLocation{}, err
		}
	}
	return p, nil
}

func (p *UserLocation) UpdateAUserLocation(db *gorm.DB) (*UserLocation, error) {

	var err error

	err = db.Debug().Model(&UserLocation{}).Where("id = ?", p.ID).Updates(UserLocation{UserLocationID: p.UserLocationID, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &UserLocation{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorID).Take(&p.Creator).Error
		if err != nil {
			return &UserLocation{}, err
		}
	}
	return p, nil
}

func (p *UserLocation) DeleteAUserLocation(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&UserLocation{}).Where("id = ? and creator_id = ?", pid, uid).Take(&UserLocation{}).Delete(&UserLocation{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("UserLocation not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
