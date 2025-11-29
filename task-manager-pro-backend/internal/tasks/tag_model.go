package tasks

type Tag struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	UserID uint   `json:"user_id" gorm:"index"`
	Name   string `json:"name" gorm:"size:50;index:idx_user_tag,unique"`
}
