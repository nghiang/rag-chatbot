package schemas

type CreateDocumentRequest struct {
	Name	string `json:"name" binding:"required"`
	FileType string `json:"file_type" binding:"required"`
}

type UpdateDocumentRequest struct {
	Name	string `json:"name"`
}