package models

type User struct {
	ID        uint64 `gorm:"primaryKey;column:id"`
	Firstname string `gorm:"column:firstname;not null"`
	Lastname  string `gorm:"column:lastname;not null"`
	Age       int    `gorm:"column:age;not null;check:age >= 0"`
}
