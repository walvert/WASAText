package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	rt.router.GET("/users", rt.wrap(rt.getUsers))

	rt.router.PUT("/users/username", rt.wrap(rt.setMyUserName))

	rt.router.PUT("/users/image", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/users/image", rt.wrap(rt.getMyPhoto))

	rt.router.GET("/uploads/:folder/images/:filename", rt.wrap(rt.getImage))

	rt.router.POST("/chats", rt.wrap(rt.createChat))
	rt.router.GET("/chats", rt.wrap(rt.getMyConversations))

	rt.router.GET("/chats/:chatId", rt.wrap(rt.getConversation))

	rt.router.PUT("/chats/:chatId/chat-name", rt.wrap(rt.setGroupName))

	rt.router.PUT("/chats/:chatId/image", rt.wrap(rt.setGroupPhoto))

	rt.router.GET("/chats/:chatId/members", rt.wrap(rt.getGroupMembers))
	rt.router.POST("/chats/:chatId/members", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/chats/:chatId/members", rt.wrap(rt.leaveGroup))

	rt.router.GET("/chats/:chatId/last-read", rt.wrap(rt.getLastRead))

	rt.router.POST("/chats/:chatId/messages", rt.wrap(rt.sendMessage))

	rt.router.DELETE("/messages/:messageId", rt.wrap(rt.deleteMessage))

	rt.router.POST("/messages/:messageId/forwards", rt.wrap(rt.forwardMessage))

	rt.router.GET("/messages/:messageId/comments", rt.wrap(rt.getComments))
	rt.router.PUT("/messages/:messageId/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/messages/:messageId/comments", rt.wrap(rt.uncommentMessage))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
