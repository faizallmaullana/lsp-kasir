package entity

import "time"

type Profiles struct {
	IdProfile string `json:"id_profile" gorm:"type:varchar(36);unique;primaryKey;not null"`
	IdUser    string `json:"id_user" gorm:"type:varchar(36);not null"`

	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	Contact  string `json:"contact" gorm:"type:varchar(120)"`
	Address  string `json:"address" gorm:"type:varchar(255)"`
	ImageUrl string `json:"photo" gorm:"type:varchar(255)"`

	IsDeleted bool      `json:"is_deleted" gorm:"type:boolean;default:false"`

	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime"`
}
