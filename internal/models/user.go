package models

// struct user represents data of user.
type User struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" binding:"required"`
	DOB         string  `json:"dob"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"createdAt"`
	Pwd         string  `json:"pwd" binding:"required"`
	Auth        string  `json:"auth"` // remove field when adding multi-login support.
}
