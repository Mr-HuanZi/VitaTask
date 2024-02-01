package dto

type ProjectListQuery struct {
	PagingQuery
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Time    []string `json:"time"`
	Deleted bool     `json:"deleted"`
}

type ProjectSimpleList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateProjectForm struct {
	Name   string `json:"name" binding:"required"`
	Leader uint64 `json:"leader"`
}

type EditProjectForm struct {
	ID     uint   `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Leader uint64 `json:"leader"`
}

type ProjectSingleId struct {
	ID uint `json:"id" binding:"required"`
}

type RelationLeaderForm struct {
	ProjectId uint
	UserId    uint64
}

type ProjectMemberListQuery struct {
	PagingQuery
	ProjectId int    `json:"project"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Role      int8   `json:"role"`
}

type ProjectMemberVO struct {
	ID        uint        `json:"id,omitempty"`
	ProjectId uint        `json:"projectId,omitempty"`
	UserId    uint64      `json:"userId,omitempty"`
	Role      int8        `json:"role,omitempty"`
	UserInfo  interface{} `json:"userInfo"`
	Project   interface{} `json:"project"`
	RoleName  []string    `json:"roleName"`
	Value     uint64      `json:"value"`
	Label     string      `json:"label"`
}

type ProjectMemberBind struct {
	ProjectId uint     `json:"project" binding:"required"`
	UserId    []uint64 `json:"users"`
	Role      int      `json:"role"`
}

type ProjectTransferForm struct {
	Project    uint   `json:"project" binding:"required"`
	Transferor uint64 `json:"transferor" binding:"required"` // 移交人
	Recipient  uint64 `json:"recipient" binding:"required"`  // 接收人
}
