package schemas

type CreateChatSessionRequest struct {
	Title           string `json:"title" binding:"required"`
	KnowledgeBaseID *uint  `json:"knowledge_base_id,omitempty"`
}

type UpdateChatSessionRequest struct {
	Title           string `json:"title"`
	KnowledgeBaseID *uint  `json:"knowledge_base_id,omitempty"`
}
