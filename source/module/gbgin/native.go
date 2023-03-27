package gbgin

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type (
	Engine         = gin.Engine
	Handler        = gin.HandlerFunc
	Context        = gin.Context
	ResponseWriter = gin.ResponseWriter
)

type CORSConfig = cors.Config
