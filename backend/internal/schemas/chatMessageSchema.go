package schemas

type CreateChatMessageRequest struct {
	Role    string `json:"role" binding:"required,oneof=user assistant system"`
	Message string `json:"message" binding:"required"`
}

type UpdateChatMessageRequest struct {
	Message string `json:"message"`
}
