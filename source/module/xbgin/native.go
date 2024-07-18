package xbgin

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/starryck/x-lib-go/source/utility/xbflow"
)

type (
	Flow     = xbflow.Flow
	BaseFlow = xbflow.BaseFlow
)

type (
	Engine         = gin.Engine
	RouterGroup    = gin.RouterGroup
	Handler        = gin.HandlerFunc
	Context        = gin.Context
	ResponseWriter = gin.ResponseWriter
)

type CORSConfig = cors.Config
