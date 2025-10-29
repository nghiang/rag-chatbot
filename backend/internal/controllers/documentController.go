package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"backend/internal/models"
	"backend/internal/queues"
	"backend/config"
	"backend/internal/services"
)


func CreateDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		kbParam := c.Param("knowledgeBaseId")
		kbID64, err := strconv.ParseUint(kbParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid knowledge base id"})
			return
		}
		kbID := uint(kbID64)

		// Parse multipart form
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}
		defer file.Close()

		fileType := c.PostForm("file_type")
		description := c.PostForm("description")

		// validate file type enum
		allowed := map[string]bool{"csv": true, "pdf": true, "graph": true}
		if !allowed[fileType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file_type must be one of csv, pdf, graph"})
			return
		}

		// prepare object name
		objectName := fmt.Sprintf("kb_%d/%d_%s", kbID, time.Now().UnixNano(), header.Filename)

		// load config (do not read env directly here)
		cfg := config.LoadConfig()
		bucket := cfg.MinioBucket

		// storage base path from config (LocalStoragePath)
		basePath := cfg.LocalStoragePath
		storageSvc := services.NewLocalStorage(basePath)

		// we need the size for PutObject; try to get from header
		var size int64 = 0
		if header != nil {
			size = header.Size
		}

		// ensure bucket (dir) exists and upload
		ctx := context.Background()
		if err := storageSvc.EnsureBucket(ctx, bucket); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to ensure bucket: " + err.Error()})
			return
		}

		if err := storageSvc.UploadObject(ctx, bucket, objectName, file, size, header.Header.Get("Content-Type")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file: " + err.Error()})
			return
		}

		// create document record with processing status
		now := time.Now()
		doc := models.Document{
			KnowledgeBaseID: kbID,
			UserID:          userID,
			Name:            header.Filename,
			FileType:        fileType,
			Description:     description,
			EmbeddingStatus: "processing",
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		if err := models.CreateDocument(&doc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// enqueue asynq task
		if err := queues.EnqueueProcessDocument(doc.ID, bucket, objectName); err != nil {
			// log but respond success because upload and DB succeeded
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue processing task: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, doc)
	}
}

func ListDocuments() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
	}
}

func GetDocumentByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
	}
}

func UpdateDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
	}
}

func DeleteDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
	}
}

