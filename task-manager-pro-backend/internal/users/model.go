package users

import "time"

type User struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"uniqueIndex"`
	Password string    `json:"-"`
	CreateAt time.Time `json:"creat_at"`
}
