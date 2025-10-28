package controllers

import (
	"time"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"backend/internal/models"
	"backend/internal/schemas"
)

func CreateKnowledgeBase() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var kbReq schemas.CreateKnowledgeBaseRequest
		if err := c.ShouldBindJSON(&kbReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		kb := models.KnowledgeBase{
			UserID:      userID,
			Name:        kbReq.Name,
			Description: kbReq.Description,
		}
		kb.CreatedAt = time.Now()
		kb.UpdatedAt = time.Now()
		if err := models.CreateKnowledgeBase(&kb); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Knowledge base created for user ID:", userID)
		c.JSON(http.StatusCreated, kb)
	}
}

func ListKnowledgeBases() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		kbs, err := models.ListKnowledgeBasesByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, kbs)
	}
}

func GetKnowledgeBaseByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid knowledge base ID"})
			return
		}
		kb, err := models.GetKnowledgeBaseByID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if kb == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "knowledge base not found"})
			return
		}
		c.JSON(http.StatusOK, kb)
	}
}

func UpdateKnowledgeBase() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid knowledge base ID"})
			return
		}

		kb, err := models.GetKnowledgeBaseByID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if kb == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "knowledge base not found"})
			return
		}
		if kb.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
			return
		}

		var kbReq schemas.UpdateKnowledgeBaseRequest
		if err := c.ShouldBindJSON(&kbReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		kb.Name = kbReq.Name
		kb.Description = kbReq.Description
		kb.UpdatedAt = time.Now()

		if err := models.UpdateKnowledgeBase(kb); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, kb)
	}
}


func DeleteKnowledgeBase() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid knowledge base ID"})
			return
		}
		if err := models.DeleteKnowledgeBase(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}