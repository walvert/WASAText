package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.POST("/session", rt.doLogin)
	rt.router.GET("/users/:userId", rt.getUser)
	rt.router.PUT("/users/:userId", rt.AuthMiddleware(rt.setMyUsername))
	rt.router.PUT("/users/:userId/image", rt.AuthMiddleware(rt.setMyPhoto))
	rt.router.POST("/users/:userId/chats", rt.AuthMiddleware(rt.sendFirstMessage))
	rt.router.POST("/users/:userId/chats/:chatId/messages", rt.AuthMiddleware(rt.sendMessage))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
