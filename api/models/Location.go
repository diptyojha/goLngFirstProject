package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Location struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Loc_Name    string    `gorm:"size:255;not null;unique" json:"Loc_Name"`
	Address     string    `gorm:"size:255;not null;" json:"Address"`
	Pincode     string    `gorm:"size:255;not null;" json:"Pincode"`
	TimeZone    string    `gorm:"size:255;not null;" json:"TimeZone"`
	Description string    `gorm:"size:255;" json:"description"`
	Creator     User      `json:"creator"`
	CreatorID   uint32    `gorm:"not null" json:"creator_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Location) Prepare() {
	p.ID = 0
	p.Loc_Name = html.EscapeString(strings.TrimSpace(p.Loc_Name))
	p.Address = html.EscapeString(strings.TrimSpace(p.Address))
	p.Pincode = html.EscapeString(strings.TrimSpace(p.Pincode))
	p.TimeZone = html.EscapeString(strings.TrimSpace(p.TimeZone))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Creator = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Location) Validate() error {

	if p.Loc_Name == "" {
		return errors.New("Required Location Name")
	}
	if p.Address == "" {
		return errors.New("Required Address")
	}
	if p.Pincode == "" {
		return errors.New("Required Pincode")
	}
	if p.TimeZone == "" {
		return errors.New("Required TimeZone")
	}
	if p.CreatorID < 1 {
		return errors.New("Required CreatorID")
	}
	return nil
}

func (p *Location) SaveLocation(db *gorm.DB) (*Location, error) {
	var err error
	err = db.Debug().Model(&Location{}).Create(&p).Error
	if err != nil {
		return &Location{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorID).Take(&p.Creator).Error
		if err != nil {
			return &Location{}, err
		}
	}
	return p, nil
}

func (p *Location) FindAllLocations(db *gorm.DB) (*[]Location, error) {
	var err error
	locations := []Location{}
	err = db.Debug().Model(&Location{}).Limit(100).Find(&locations).Error
	if err != nil {
		return &[]Location{}, err
	}
	if len(locations) > 0 {
		for i, _ := range locations {
			err := db.Debug().Model(&User{}).Where("id = ?", locations[i].CreatorID).Take(&locations[i].Creator).Error
			if err != nil {
				return &[]Location{}, err
			}
		}
	}
	return &locations, nil
}

func (p *Location) FindLocationByID(db *gorm.DB, pid uint64) (*Location, error) {
	var err error
	err = db.Debug().Model(&Location{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Location{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorID).Take(&p.Creator).Error
		if err != nil {
			return &Location{}, err
		}
	}
	return p, nil
}

func (p *Location) UpdateALocation(db *gorm.DB) (*Location, error) {

	var err error

	err = db.Debug().Model(&Location{}).Where("id = ?", p.ID).Updates(Location{Loc_Name: p.Loc_Name, Address: p.Address, Pincode: p.Pincode, TimeZone: p.TimeZone, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Location{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorID).Take(&p.Creator).Error
		if err != nil {
			return &Location{}, err
		}
	}
	return p, nil
}

func (p *Location) DeleteALocation(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Location{}).Where("id = ? and creator_id = ?", pid, uid).Take(&Location{}).Delete(&Location{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Location not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
