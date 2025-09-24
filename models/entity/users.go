package entity

import "time"

type Users struct {
	IdUser   string `json:"id_user" gorm:"type:varchar(36);unique;primaryKey;not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
	Role     string `json:"role" gorm:"type:varchar(50);not null"`

	Transactions []Transactions `json:"transactions" gorm:"foreignKey:IdUser;references:IdUser;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Sessions     []Sessions     `json:"sessions" gorm:"foreignKey:IdUser;references:IdUser;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Profiles     []Profiles     `json:"profiles" gorm:"foreignKey:IdUser;references:IdUser;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	IsDeleted bool `json:"is_deleted" gorm:"type:boolean;default:false"`

	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime"`
}
