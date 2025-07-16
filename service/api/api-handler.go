package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/users", rt.wrap(rt.getUsers))
	rt.router.PUT("/users", rt.setMyUsername)
	rt.router.PUT("/users/image", rt.setMyPhoto)
	rt.router.POST("/users/media", rt.authMiddleware(rt.uploadMessageMedia))
	rt.router.GET("/uploads/:folder/:filename", rt.wrap(rt.getImage))
	rt.router.POST("/chats", rt.createChat)
	rt.router.GET("/chats", rt.getMyConversations)
	rt.router.GET("/chats/:chatId", rt.getConversation)
	rt.router.PUT("/chats/:chatId", rt.setGroupName)
	rt.router.PUT("/chats/:chatId/image", rt.setGroupPhoto)
	rt.router.POST("/chats/:chatId/members", rt.addToGroup)
	rt.router.DELETE("/chats/:chatId/members", rt.leaveGroup)
	rt.router.GET("/chats/:chatId/last-read", rt.getLastRead)
	rt.router.POST("/chats/:chatId/messages", rt.wrap(rt.sendMessage))
	rt.router.DELETE("/chats/:chatId/messages/:messageId", rt.authMiddleware(rt.authDeleteMessage(rt.deleteMessage)))
	rt.router.POST("/chats/:chatId/messages/:messageId", rt.authMiddleware(rt.forwardMessage))
	rt.router.DELETE("/chats/:chatId/messages/:messageId/comments", rt.authMiddleware(rt.deleteComment))
	rt.router.PUT("/chats/:chatId/messages/:messageId/comments", rt.authMiddleware(rt.commentMessage))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
