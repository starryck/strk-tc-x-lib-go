package xbgin

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

type RouterStem struct {
	Path     string
	Handlers []Handler
	Leaves   []RouterLeaf
	Stems    []RouterStem
}

type RouterLeaf struct {
	Method   string
	Path     string
	Handlers []Handler
}

func (router *Router) GetEngine() *Engine {
	return router.engine
}

func (router *Router) GetCORSConfig() *CORSConfig {
	return router.corsConfig
}

func (router *Router) UseMiddlewares() {
	router.SetMiddlewares(router.NewMiddlewares()...)
}

func (router *Router) SetMiddlewares(handlers ...Handler) {
	router.engine.Use(handlers...)
}

func (router *Router) NewMiddlewares() []Handler {
	return []Handler{
		gin.Recovery(),
		cors.New(*router.corsConfig),
		GraceMiddleware,
		RecordMiddleware,
		ResponseMiddleware,
	}
}

func (router *Router) SetRouterGroup(stems ...RouterStem) {
	router.setRouterGroup(&router.engine.RouterGroup, stems...)
}

func (router *Router) setRouterGroup(group *RouterGroup, stems ...RouterStem) {
	for _, stem := range stems {
		subgroup := group.Group(stem.Path, stem.Handlers...)
		for _, leaf := range stem.Leaves {
			subgroup.Handle(leaf.Method, leaf.Path, leaf.Handlers...)
		}
		router.setRouterGroup(subgroup, stem.Stems...)
	}
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
	engine.NoRoute(NoRouteHandler)
	builder.router.engine = engine
	return builder
}

func (builder *routerBuilder) setCORSConfig() *routerBuilder {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	builder.router.corsConfig = &corsConfig
	return builder
}
