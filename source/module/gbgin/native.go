package gbgin

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/forbot161602/pbc-golang-lib/source/utility/gbflow"
)

type (
	Flow     = gbflow.Flow
	BaseFlow = gbflow.BaseFlow
)

type (
	Engine         = gin.Engine
	Handler        = gin.HandlerFunc
	Context        = gin.Context
	ResponseWriter = gin.ResponseWriter
)

type CORSConfig = cors.Config
