package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.POST("/session", rt.login)
	rt.router.GET("/users/:id", rt.getUser)
	rt.router.PUT("/users/:id", rt.editUsername)
	rt.router.PUT("/users/:id/image", rt.uploadUserImage)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
