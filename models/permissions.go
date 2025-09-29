package models

type Permission struct {
	ID          int     `gorm:"primaryKey;autoIncrement;column:id"`
	Code        string  `gorm:"type:varchar;uniqueIndex;not null;column:code"`
	Description *string `gorm:"type:text;column:description"`
}
