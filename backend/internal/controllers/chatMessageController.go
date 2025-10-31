package controllers

import (
	"backend/internal/models"
	"backend/internal/schemas"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateChatMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		sessionIDParam := c.Param("id")
		sessionID, err := strconv.ParseUint(sessionIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
			return
		}

		// Verify that the session exists and belongs to the user
		session, err := models.GetChatSessionByID(uint(sessionID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if session == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "chat session not found"})
			return
		}
		if session.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
			return
		}

		var messageReq schemas.CreateChatMessageRequest
		if err := c.ShouldBindJSON(&messageReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		message := models.ChatMessage{
			SessionID: uint(sessionID),
			Role:      messageReq.Role,
			Message:   messageReq.Message,
			CreatedAt: time.Now(),
		}

		if err := models.CreateChatMessage(&message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Update session's UpdatedAt timestamp
		session.UpdatedAt = time.Now()
		if err := models.UpdateChatSession(session); err != nil {
			// Log the error but don't fail the request
		}

		c.JSON(http.StatusCreated, message)
	}
}

func ListChatMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		sessionIDParam := c.Param("id")
		sessionID, err := strconv.ParseUint(sessionIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
			return
		}

		// Verify that the session exists and belongs to the user
		session, err := models.GetChatSessionByID(uint(sessionID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if session == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "chat session not found"})
			return
		}
		if session.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
			return
		}

		messages, err := models.ListChatMessagesBySessionID(uint(sessionID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, messages)
	}
}

func GetChatMessageByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		sessionIDParam := c.Param("id")
		sessionID, err := strconv.ParseUint(sessionIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
			return
		}

		idParam := c.Param("messageId")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
			return
		}

		// Verify session ownership
		session, err := models.GetChatSessionByID(uint(sessionID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if session == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "chat session not found"})
			return
		}
		if session.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
			return
		}

		message, err := models.GetChatMessageByID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if message == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
			return
		}

		// Verify message belongs to the session
		if message.SessionID != uint(sessionID) {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found in this session"})
			return
		}

		c.JSON(http.StatusOK, message)
	}
}

func UpdateChatMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		sessionIDParam := c.Param("id")
		sessionID, err := strconv.ParseUint(sessionIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
			return
		}

		idParam := c.Param("messageId")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
			return
		}

		// Verify session ownership
		session, err := models.GetChatSessionByID(uint(sessionID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if session == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "chat session not found"})
			return
		}
		if session.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
			return
		}

		message, err := models.GetChatMessageByID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if message == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
			return
		}

		// Verify message belongs to the session
		if message.SessionID != uint(sessionID) {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found in this session"})
			return
		}

		var messageReq schemas.UpdateChatMessageRequest
		if err := c.ShouldBindJSON(&messageReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		message.Message = messageReq.Message

		if err := models.UpdateChatMessage(message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, message)
	}
}

func DeleteChatMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		sessionIDParam := c.Param("id")
		sessionID, err := strconv.ParseUint(sessionIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
			return
		}

		idParam := c.Param("messageId")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
			return
		}

		// Verify session ownership
		session, err := models.GetChatSessionByID(uint(sessionID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if session == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "chat session not found"})
			return
		}
		if session.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
			return
		}

		message, err := models.GetChatMessageByID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if message == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
			return
		}

		// Verify message belongs to the session
		if message.SessionID != uint(sessionID) {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found in this session"})
			return
		}

		if err := models.DeleteChatMessage(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
