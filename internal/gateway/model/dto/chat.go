package dto

type ChatSendUserForm struct {
	Userid string `json:"userid" binding:"required"`
	Msg    string `json:"msg" binding:"required"`
}

type ChatSendUsersForm struct {
	Users []string `json:"users" binding:"required"`
	Msg   string   `json:"msg" binding:"required"`
}
