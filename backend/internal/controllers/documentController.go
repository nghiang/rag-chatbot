package controllers

import (
	"backend/config"
	"backend/internal/models"
	"backend/internal/queues"
	"backend/internal/services"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
		allowed := map[string]bool{"csv": true, "doc": true, "graph": true}
		if !allowed[fileType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file_type must be one of csv, doc, graph"})
			return
		}

		// prepare object name
		objectName := fmt.Sprintf("kb_%d/%s", kbID, header.Filename)

		// load config
		cfg := config.LoadConfig()
		bucket := cfg.MinioBucket

		// Use MinIO storage
		storageSvc, err := services.NewMinIOStorage(
			cfg.MinioEndpoint,
			cfg.MinioAccessKey,
			cfg.MinioSecretKey,
			cfg.MinioUseSSL,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize MinIO storage: " + err.Error()})
			return
		}

		// ensure bucket exists and upload
		ctx := context.Background()
		if err := storageSvc.EnsureBucket(ctx, bucket); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to ensure bucket: " + err.Error()})
			return
		}

		contentType := header.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		if err := storageSvc.UploadObject(ctx, bucket, objectName, file, header.Size, contentType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file: " + err.Error()})
			return
		}

		log.Printf("File uploaded to MinIO: bucket=%s, object=%s, size=%d", bucket, objectName, header.Size)

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
		if err := queues.EnqueueProcessDocument(kbID, doc.ID, description, bucket, objectName, fileType); err != nil {
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
