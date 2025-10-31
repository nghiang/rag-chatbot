package controllers

import (
	"backend/internal/models"
	"backend/internal/schemas"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateChatSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var sessionReq schemas.CreateChatSessionRequest
		if err := c.ShouldBindJSON(&sessionReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		session := models.ChatSession{
			UserID:          userID,
			Title:           sessionReq.Title,
			KnowledgeBaseID: sessionReq.KnowledgeBaseID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := models.CreateChatSession(&session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, session)
	}
}

func ListChatSessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		sessions, err := models.ListChatSessionsByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, sessions)
	}
}

func GetChatSessionByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat session ID"})
			return
		}

		session, err := models.GetChatSessionByID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if session == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "chat session not found"})
			return
		}

		// Check if the session belongs to the user
		if session.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func UpdateChatSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat session ID"})
			return
		}

		session, err := models.GetChatSessionByID(uint(id))
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

		var sessionReq schemas.UpdateChatSessionRequest
		if err := c.ShouldBindJSON(&sessionReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		session.Title = sessionReq.Title
		session.KnowledgeBaseID = sessionReq.KnowledgeBaseID
		session.UpdatedAt = time.Now()

		if err := models.UpdateChatSession(session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func DeleteChatSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat session ID"})
			return
		}

		session, err := models.GetChatSessionByID(uint(id))
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

		if err := models.DeleteChatSession(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
