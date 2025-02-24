package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.doLogin)
	rt.router.PUT("/users/:userId", rt.authMiddleware(rt.setMyUsername))
	rt.router.PUT("/users/:userId/image", rt.authMiddleware(rt.setMyPhoto))
	rt.router.POST("/users/:userId/chats", rt.authMiddleware(rt.sendFirstMessage))
	rt.router.GET("/users/:userId/chats", rt.authMiddleware(rt.getMyConversations))
	rt.router.GET("/users/:userId/chats/:chatId", rt.authMiddleware(rt.getConversation))
	rt.router.PUT("/users/:userId/chats/:chatId", rt.authMiddleware(rt.setGroupName))
	rt.router.PUT("/users/:userId/chats/:chatId/image", rt.authMiddleware(rt.setGroupPhoto))
	rt.router.POST("/users/:userId/chats/:chatId/members", rt.authMiddleware(rt.addToGroup))
	rt.router.DELETE("/users/:userId/chats/:chatId/members", rt.authMiddleware(rt.LeaveGroup))
	rt.router.POST("/users/:userId/chats/:chatId/messages", rt.authMiddleware(rt.sendMessage))
	rt.router.DELETE("/users/:userId/chats/:chatId/messages/:messageId", rt.authMiddleware(rt.authDeleteMessage(rt.deleteMessage)))
	rt.router.POST("/users/:userId/chats/:chatId/messages/:messageId", rt.authMiddleware(rt.forwardMessage))
	rt.router.DELETE("/users/:userId/chats/:chatId/messages/:messageId/comments", rt.authMiddleware(rt.deleteComment))
	rt.router.PUT("/users/:userId/chats/:chatId/messages/:messageId/comments", rt.authMiddleware(rt.commentMessage))

	/*
		 doLogin (see simplified login) √
		• setMyUserName 				√
		• getMyConversations			√
		• getConversation				√
		• sendMessage					√
		• forwardMessage				√
		• commentMessage				√
		• uncommentMessage				√
		• deleteMessage					√
		• addToGroup					√
		• leaveGroup					√
		• setGroupName					√
		• setMyPhoto					√
		• setGroupPhoto					√
	*/

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
