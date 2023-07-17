package models

import (
	"database/sql/driver"
	"encoding/json"
	// "errors"
	"fmt"
)

type TableArray []Table
type Tab struct {
	ID         uint       `json:"no" gorm:"primary_key"`
	Tables     TableArray `json:"tables" gorm:"type:text"`
	QuestionID uint       `json:"-"`
	Question   *Question  `json:"question,omitempty" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

func (sla TableArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

// Scan Unmarshal
func (a *TableArray) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %#v", value)
	}
	return json.Unmarshal(bytes, &a)
}
