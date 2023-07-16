package models
// ! this is not gonna be in the database
// ! it's part of the tab model
type Table struct {
	Title string     `json:"tableTitle"`
	Rows  [][]string `json:"rows"`
}
