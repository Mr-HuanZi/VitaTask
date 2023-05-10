package types

type SimpleMemberList struct {
	Id           uint64 `json:"value,omitempty"`
	UserLogin    string `json:"text,omitempty"`
	UserNickname string `json:"label,omitempty"`
	Avatar       string `json:"avatar"`
}

type MemberCreate struct {
	Nickname string `json:"nickname,omitempty" binding:"required"`
	Username string `json:"username,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
	Email    string `json:"email,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
}

type MemberListsQuery struct {
	PagingQuery
	Id       uint64 `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Status   int    `json:"status,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	Email    string `json:"email,omitempty"`
}
