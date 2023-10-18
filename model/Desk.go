package model

import "github.com/jinzhu/gorm"

type Desk struct {
	gorm.Model
	Name string
}
