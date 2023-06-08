package models

import (
	U "docker/utils"
	"time"
)

var FileTypes = map[string]uint16{
	"plan":        1,
	"certificate": 2,
	"attachment":  3,
}
var IntFileTypes = map[uint16]string{
	1: "plan",
	2: "certificate",
	3: "attachment",
}

type File struct {
	ID   uint   `gorm:"primaryKey"`
	Type uint16 `json:"type"` // like attachment, certificate and etc...
	Name string `json:"name"`
	// relations
	LawID     uint
	Law       Law        `json:"law"`
	CreatedAt *time.Time `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"not null;default:now()"`
}
type FileMinimal struct {
	ID uint `json:"id"`
	// Type string `json:"type"` // like attachment, certificate and etc...
	Type uint16 `json:"type"` // like attachment, certificate and etc...
	URL  string `json:"name"`
	// relations
	CreatedAt *time.Time `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"not null;default:now()"`
}

// tips: if convert is array, no need to call be reference
// but if convert is about one object, it would be better if we call them as
// -- call by refrence
func FileToFileMinimal(files []File) []FileMinimal {
	var minimalFiles []FileMinimal
	for i := 0; i < len(files); i++ {
		minimalFile := FileMinimal{
			ID: files[i].ID,
			// Type:      IntFileTypes[files[i].Type],
			Type:      files[i].Type,
			URL:       U.BaseURL + "/public/uploads/" + files[i].Name,
			CreatedAt: files[i].CreatedAt,
			UpdatedAt: files[i].UpdatedAt,
		}
		minimalFiles = append(minimalFiles, minimalFile)
	}
	return minimalFiles
}
