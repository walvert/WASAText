package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.POST("/session", rt.doLogin)
	rt.router.GET("/users/:id", rt.getUser)
	rt.router.PUT("/users/:id", rt.AuthMiddleware(rt.setMyUsername))
	rt.router.PUT("/users/:id/image", rt.AuthMiddleware(rt.setMyPhoto))
	rt.router.POST("/users/:id/chats", rt.createChat)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
