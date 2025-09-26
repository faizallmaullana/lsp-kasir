package entity

import "time"

type Images struct {
	IdImage     string    `json:"id_image" gorm:"type:varchar(36);unique;primaryKey;not null"`
	FileName    string    `json:"file_name" gorm:"type:varchar(255)"`
	ContentType string    `json:"content_type" gorm:"type:varchar(120)"`
	Size        int64     `json:"size"`
	Data        []byte    `json:"-" gorm:"type:bytea"`
	IsDeleted   bool      `json:"is_deleted" gorm:"type:boolean;default:false"`
	Timestamp   time.Time `json:"timestamp" gorm:"autoCreateTime"`
}
