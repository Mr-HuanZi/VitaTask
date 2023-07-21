package response

const (
	SystemFail             = 0   // 系统错误
	FormVerificationFailed = 101 // 表单校验失败
	SignatureMissing       = 102 // 签名丢失
	DbQueryError           = 103 // 数据库查询错误
	DbExecuteError         = 104 // 数据库操作执行错误
	NotLoggedIn            = 105 // 未登录

	LoginSingGenerateFail = 201 // 签名生成失败
	LoginPassError        = 202 // 用户名或密码不正确
	RegUsernameExists     = 203 // 用户名已存在
	RegFail               = 204 // 注册失败
	RegPassFormatError    = 205 // 密码格式错误
	EmptyUsernameOrPass   = 206 // 空的用户名或密码
	PassError             = 207 // 密码错误
	NotInputtedMobile     = 208 // 未输入手机号
	NotInputtedEmail      = 209 // 未输入电子邮箱地址

	UserNotFound        = 1000 // 用户不存在
	MemberCreateFail    = 1001 // 成员创建失败
	UserDisabled        = 1002 // 用户已禁用
	AvatarNotUploaded   = 1003 // 头像未上传
	CurrUserNotSuper    = 1004 // 当前用户不是超级用户
	UserSuperChangeSelf = 1005 // 不能改变自己的超级用户状态

	ProjectCreateFail            = 2001 // 项目创建失败
	ProjectUpdateFail            = 2002 // 项目更新失败
	ProjectStared                = 2003 // 项目已经Star
	ProjectLeaderNotExist        = 2004 // 项目负责人不存在
	ProjectNotExist              = 2005 // 项目不存在
	ProjectDeleteFail            = 2006 // 项目删除失败
	ProjectArchived              = 2007 // 项目已归档
	ProjectArchiveFail           = 2008 // 项目归档失败
	ProjectMemberQueryFail       = 2009 // 项目成员查询失败
	ProjectNotArchived           = 2010 // 项目未归档
	ProjectUnArchiveFail         = 2011 // 项目取消归档失败
	ProjectMultipleSpecialMember = 2012 // 多个负责人或创建人
	ProjectRoleNonExistent       = 2013 // 项目角色不存在
	ProjectLeaderRemove          = 2014 // 移除项目负责人

	TaskCreateFail            = 2100 // 任务创建失败
	TaskStatusNotExist        = 2101 // 任务状态不存在
	TaskRoleNonExistent       = 2102 // 任务角色不存在
	TaskMultipleSpecialMember = 2103 // 多个负责人或创建人
	TaskNotExist              = 2104 // 任务不存在
	TaskCreatorRemove         = 2105 // 移除创建人
	TaskDeleteFail            = 2106 // 任务删除失败
	TaskStatusProcessing      = 2107 // 任务仍在进行中

	TaskGroupNotExist = 2200 // 任务组不存在

	TaskOperatorTypeIllegal = 2300 // 非法的任务操作类型

	MemberNotInProject     = 3000 // 成员不在项目内
	MemberNotProjectLeader = 3001 // 成员不是项目负责人

	FilesLimitExceeded = 4000 // 文件大小超出限制
	FilesSuffixError   = 4001 // 文件后缀错误

	DialogNotExist        = 5000 // 对话不存在
	NotInDialog           = 5001 // 不是该对话成员
	DialogTypeError       = 5002 // 对话类型错误
	DialogMemberEmpty     = 5003 // 对话成员为空
	DialogMemberIsMe      = 5004 // 不能和自己对话
	DialogCreateFail      = 5005 // 对话创建失败
	DialogC2COvercrowding = 5006 // C2C对话超员
	IsInDialog            = 5007 // 成员已在对话中
	JoinDialogFail        = 5008 // 加入对话失败
	DialogKeep1Member     = 5009 // 对话至少保留一个成员
	DialogDeleteFail      = 5010 // 对话删除失败

	TimeParseFail            = 9000 // 时间解析失败
	ElementQuantityTooLittle = 9001 // 元素数量太少
	ElementQuantityTooMany   = 9002 // 元素数量太多
	StartTimeGtEndTime       = 9003 // 开始时间大于结束时间
	TimeSpanTooLong          = 9004 // 时间跨度太长
	TooFewElements           = 9005 // 元素太少。一般用于切片或map元素个数判断
)

var codeMap = map[int]string{
	SystemFail:             "系统错误",
	FormVerificationFailed: "表单校验失败",
	SignatureMissing:       "签名失效",
	DbQueryError:           "数据库查询错误",
	DbExecuteError:         "数据库操作执行错误",
	NotLoggedIn:            "用户未登录",

	LoginSingGenerateFail: "签名生成失败",
	LoginPassError:        "用户名或密码不正确",
	RegUsernameExists:     "用户名已存在",
	RegFail:               "注册失败",
	RegPassFormatError:    "密码必须包含大小写字母和数字的组合，可以使用特殊字符，长度在8-16之间",
	EmptyUsernameOrPass:   "用户名或密码不能为空",
	PassError:             "密码错误",
	NotInputtedMobile:     "未输入手机号",
	NotInputtedEmail:      "未输入电子邮箱地址",

	UserNotFound:        "用户不存在",
	MemberCreateFail:    "成员创建失败",
	UserDisabled:        "用户已禁用",
	AvatarNotUploaded:   "头像未上传",
	CurrUserNotSuper:    "您不是超级用户",
	UserSuperChangeSelf: "不能改变自己的超级用户状态",

	ProjectCreateFail:            "项目创建失败",
	ProjectUpdateFail:            "项目更新失败",
	ProjectStared:                "项目已经Star",
	ProjectLeaderNotExist:        "项目负责人不存在",
	ProjectNotExist:              "项目不存在",
	ProjectDeleteFail:            "项目删除失败",
	ProjectArchived:              "项目已是归档状态",
	ProjectArchiveFail:           "项目归档失败",
	ProjectMemberQueryFail:       "项目成员查询失败",
	ProjectNotArchived:           "项目不是归档状态",
	ProjectUnArchiveFail:         "取消项目归档失败",
	ProjectMultipleSpecialMember: "一个项目只能有一个负责人或创建人",
	ProjectRoleNonExistent:       "项目角色不存在",
	ProjectLeaderRemove:          "不得移除项目负责人",

	TaskCreateFail:            "项目创建失败",
	TaskNotExist:              "任务不存在",
	TaskRoleNonExistent:       "任务角色不存在",
	TaskMultipleSpecialMember: "一个任务只能有一个负责人或创建人",
	TaskStatusNotExist:        "任务状态不存在",
	TaskCreatorRemove:         "不得移除任务创建人",
	TaskDeleteFail:            "任务删除失败",
	TaskStatusProcessing:      "任务仍在进行中",

	TaskGroupNotExist: "任务组不存在",

	TaskOperatorTypeIllegal: "非法的任务操作类型",

	MemberNotInProject:     "成员不在项目内",
	MemberNotProjectLeader: "成员不是项目负责人",

	FilesLimitExceeded: "文件大小超出限制",
	FilesSuffixError:   "文件后缀错误",

	DialogNotExist:        "对话不存在",
	NotInDialog:           "不是该对话成员",
	DialogTypeError:       "对话类型错误",
	DialogMemberEmpty:     "对话成员为空",
	DialogMemberIsMe:      "不能和自己对话",
	DialogCreateFail:      "对话创建失败",
	DialogC2COvercrowding: "C2C对话只允许一个聊天对象",
	IsInDialog:            "成员已在对话中",
	JoinDialogFail:        "加入对话失败",
	DialogKeep1Member:     "对话至少保留一个成员",
	DialogDeleteFail:      "对话删除失败",

	TimeParseFail:            "时间解析失败",
	ElementQuantityTooLittle: "元素数量太少",
	ElementQuantityTooMany:   "元素数量太多",
	StartTimeGtEndTime:       "开始时间大于结束时间",
	TimeSpanTooLong:          "时间跨度太长",
	TooFewElements:           "元素太少",
}

// GetMessage 根据状态码获取消息
func GetMessage(code int) string {
	if msg, ok := codeMap[code]; !ok {
		return "未知异常"
	} else {
		return msg
	}
}
