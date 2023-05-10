package types

type DialogSendTextDto struct {
	DialogId uint   `json:"dialog_id" binding:"required"`
	Content  string `json:"content" binding:"required,gt=0"`
}

type DialogCreateDto struct {
	Name    string   `json:"name" binding:"required"`
	Type    string   `json:"type" binding:"required"`
	Members []uint64 `json:"members" binding:"required,gt=0"`
}

type DialogIdDto struct {
	DialogId uint `json:"dialog_id" binding:"required"`
}
