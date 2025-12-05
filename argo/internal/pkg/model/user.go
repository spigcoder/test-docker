package model

import "time"

// BaseModel 为基础模型，每个表都需要包含以下字段：
// - ID：每个记录的唯一标识符（主键）
// - CreatedAt：创建时间，在创建记录时自动设置为当前时间
// - UpdatedAt：更新时间，每当记录更新时都会更新为当前时间
type BaseModel struct {
	ID        uint      `gorm:"primaryKey;column:id;"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime"`
}

type Profile struct {
	BaseModel

	Name   string `gorm:"column:name"`
	Avatar string `gorm:"column:avatar"`

	UserID uint `gorm:"column:userId;unique"`
}
