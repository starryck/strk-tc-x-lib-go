package xbgorm

import "gorm.io/gorm"

type (
	Client    = gorm.DB
	Config    = gorm.Config
	Dialector = gorm.Dialector
	Session   = gorm.Session
	Model     = gorm.Model
	DeleteAt  = gorm.DeletedAt
)
