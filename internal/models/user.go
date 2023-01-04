package models

// struct user represents data of user.
type User struct {
	ID          string     `json:"id"`
	Name        string     `json:"name" binding:"required"`
	DOB         string     `json:"dob"`
	Address     string     `json:"address"` //change to arr of strings
	Coordinates [2]float64 `json:"coordinates" binding:"required"`
	Description string     `json:"description"`
	CreatedAt   string     `json:"createdAt"`
	Pwd         string     `json:"pwd" binding:"required"`
	Auth        string     `json:"auth"` // remove field when adding multi-login support.
}

// struct UpdatableUser represents only those data of users that are updatable.
type UpdatableUser struct {
	Description string     `json:"description"`
	Address     string     `json:"address"`
	Coordinates [2]float64 `json:"coordinates"`
	Pwd         string     `json:"pwd"`
}

/* // pass *struct as parameter s
func StructIterator(s any) {

} */

// add bindJSON as receiver method to above structs.
