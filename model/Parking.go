package model

import (
	"database/sql/driver"

	"github.com/jinzhu/gorm"
)

type resourceType string

const (
	PARKING resourceType = "PARKING"
	DESK    resourceType = "DESK"
)

func (ct *resourceType) Scan(value interface{}) error {
	*ct = resourceType(value.([]byte))
	return nil
}

func (ct resourceType) Value() (driver.Value, error) {
	return string(ct), nil
}

type Resource struct {
	gorm.Model
	Name         string
	ResourceType resourceType `gorm:"type:resource_type"`
}
