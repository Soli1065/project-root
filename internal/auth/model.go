// internal/auth/model.go
package model

import "github.com/jinzhu/gorm"

type User struct {
    gorm.Model
    Username string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
}
