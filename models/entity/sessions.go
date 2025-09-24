package entity

import "time"

type Sessions struct {
	IdSession string `json:"id_session" gorm:"type:varchar(36);unique;primaryKey;not null"`
	IdUser    string `json:"id_user" gorm:"type:varchar(36);not null"`

	IsLogedIn bool `json:"is_loged_in"`

	IsDeleted bool `json:"is_deleted" gorm:"type:boolean;default:false"`

	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime"`
}
