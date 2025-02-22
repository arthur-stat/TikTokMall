package model

import "time"

type Base struct {
    ID uint32 `gorm:"primaryKey;autoIncrement;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoCreateTime"`
    DeletedAt time.Time `gorm:"autoCreateTime"`
}
