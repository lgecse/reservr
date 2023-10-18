package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Reservation struct {
	gorm.Model
	Resource Resource
	Date     time.Time
	User     User
}
