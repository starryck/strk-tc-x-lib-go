package xbgorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/starryck/x-lib-go/source/core/utility/xberror"
)

type (
	Client      = gorm.DB
	Config      = gorm.Config
	Dialector   = gorm.Dialector
	Session     = gorm.Session
	Model       = gorm.Model
	DeleteAt    = gorm.DeletedAt
	Association = gorm.Association
	Expression  = clause.Expression
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

func IsErrRecordNotFound(err error) bool {
	return xberror.Is(err, ErrRecordNotFound)
}
