package models

import (
	"database/sql/driver"
	"encoding/json"
)

type TabArray []Tab
type Tab struct {
	ID         uint      `json:"no" gorm:"primary_key"`
	Tables     TabArray  `json:"tables" gorm:"type:text"`
	QuestionID uint      `json:"-"`
	Question   *Question `json:"question,omitempty" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

func (sla *TabArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla TabArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}
