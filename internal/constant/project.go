package constant

const (
	ProjectCreate = 1 << iota // 项目创建者	0001 -> 1
	ProjectLeader             // 项目负责人	0010 -> 2
	ProjectMember             // 项目成员 0100 -> 4
	ProjectStar               // 项目收藏者 1000 -> 8
)

const (
	TaskCreator = 1 << iota // 任务创建者 00001 -> 1
	TaskLeader              // 任务负责人 00010 -> 2
	TaskMember              // 任务普通成员 00100 -> 4
	TaskFollow              // 任务关注者 01000 -> 8
	TaskTester              // 任务测试员 10000 -> 16
)

const (
	ProjectNotArchive = iota // 项目未归档
	ProjectArchived          // 项目已归档
)

const (
	TaskStatusProcessing = iota
	TaskStatusCompleted
	TaskStatusArchived
)

const (
	TaskOperatorCreate             = "create"
	TaskOperatorUpdate             = "update"
	TaskOperatorDelete             = "delete"
	TaskOperatorStatus             = "status"
	TaskOperatorRemoveLeader       = "remove_leader"
	TaskOperatorChangeLeader       = "change_leader"
	TaskOperatorChangeCollaborator = "change_collaborator"
)

var projectRole = map[int]string{
	ProjectCreate: "创建人",
	ProjectLeader: "负责人",
	ProjectMember: "成员",
	ProjectStar:   "收藏者",
}

var taskRole = map[int]string{
	TaskCreator: "创建人",
	TaskLeader:  "负责人",
	TaskMember:  "成员",
	TaskFollow:  "关注者",
	TaskTester:  "测试员",
}

var taskStatus = map[int]string{
	TaskStatusProcessing: "进行中",
	TaskStatusCompleted:  "已完成",
	TaskStatusArchived:   "已归档",
}

func GetProjectRoles() map[int]string {
	return projectRole
}

func GetProjectRole(role int) string {
	projectRole := GetProjectRoles()
	if item, ok := projectRole[role]; ok {
		return item
	}

	return ""
}

func GetTaskRoles() map[int]string {
	return taskRole
}

func GetTaskStatus() map[int]string {
	return taskStatus
}

func GetTaskLogOperatorMaps() map[string]string {
	return map[string]string{
		TaskOperatorCreate:             "创建任务",
		TaskOperatorUpdate:             "更新任务",
		TaskOperatorDelete:             "删除任务",
		TaskOperatorStatus:             "修改状态",
		TaskOperatorRemoveLeader:       "移除负责人",
		TaskOperatorChangeLeader:       "变更负责人",
		TaskOperatorChangeCollaborator: "变更协作人",
	}
}
