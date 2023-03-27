package gbgin

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *Router {
	router := (&routerBuilder{}).
		initialize().
		setEngine().
		setCORSConfig().
		build()
	return router
}

type Router struct {
	engine     *Engine
	corsConfig *CORSConfig
}

type rebuildHandlers = func(handlers []Handler) []Handler

func (router *Router) GetEngine() *Engine {
	return router.engine
}

func (router *Router) GetCORSConfig() *CORSConfig {
	return router.corsConfig
}

func (router *Router) SetMiddlewares(rebuild rebuildHandlers) {
	router.engine.Use(rebuild([]Handler{
		gin.Recovery(),
		cors.New(*router.corsConfig),
	})...)
}

type routerBuilder struct {
	router *Router
}

func (builder *routerBuilder) build() *Router {
	return builder.router
}

func (builder *routerBuilder) initialize() *routerBuilder {
	builder.router = &Router{}
	return builder
}

func (builder *routerBuilder) setEngine() *routerBuilder {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	builder.router.engine = engine
	return builder
}

func (builder *routerBuilder) setCORSConfig() *routerBuilder {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	builder.router.corsConfig = &corsConfig
	return builder
}
