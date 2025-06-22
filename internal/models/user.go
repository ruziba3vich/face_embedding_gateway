package models

type User struct {
	ID        string `json:"id" gorm:"primaryKey;type:uuid"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Password  string `json:"password"`
	PictureID string `json:"pic_id"`
}
