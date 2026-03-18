package model

import "time"

type TokenBlacklist struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Token     string    `gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
