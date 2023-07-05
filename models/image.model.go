package models

// polymorphic image
type Image struct {
	ID int `json:"-" gorm:"primary_key"`
	// image name, but after preloading the images, we need to add the base url to the image
	// and return the url of the image
	Name      string `json:"url"`
	OwnerID   int    `json:"-"`
	OwnerType string `json:"-"`
}
