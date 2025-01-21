package validator

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var formVerificationFailed = map[string]string{
	// 登录表单验证
	"LoginForm.Username.required": "请输入用户名",
	"LoginForm.Password.required": "请输入密码",
	"LoginForm.Code.required":     "请输入验证码",
	// 用户注册
	"UserRegisterForm.Username.required":        "请输入用户名",
	"UserRegisterForm.Password.required":        "请输入密码",
	"UserRegisterForm.Password.eqfield":         "密码与确认密码不一致",
	"UserRegisterForm.Password.gte":             "密码长度最小 5 位，最大 20 位",
	"UserRegisterForm.Password.lte":             "密码长度最小 5 位，最大 20 位",
	"UserRegisterForm.ConfirmPassword.required": "请输入确认密码",
	"UserRegisterForm.UserNickname.required":    "请输入用户昵称",
	"UserRegisterForm.UserEmail.email":          "邮箱格式不正确",
	"UserRegister.UsernameExists":               "用户名已存在",
	"UserRegister.Fail":                         "系统错误，注册失败",
	// 项目
	"CreateProjectForm.Name.required": "请填写项目名称",
	"CreateProjectForm.ID.required":   "缺少项目ID参数",
	// 任务组
	"TaskGroupForm.ProjectId.required":             "缺少项目ID参数",
	"TaskGroupQuery.PagingQuery.PageSize.required": "缺少PageSize参数",
	// 任务
	"TaskCreateForm.ProjectId.required": "缺少项目ID参数",
	// 对话
	"DiaLogSendTextDto.DialogId.required": "请选择对话",
	"DialogSendTextDto.Token.required":    "缺少WebsocketToken",
	"DiaLogSendTextDto.Content.required":  "请填写对话正文",
	"DiaLogSendTextDto.Content.gt":        "请填写对话正文",
	// DialogCreateDto 创建对话Dto
	"DialogCreateDto.Name.required":    "请填写对话名称",
	"DialogCreateDto.Type.required":    "请填写对话类型",
	"DialogCreateDto.Members.required": "请选择对话成员",
	"DialogCreateDto.Members.gt":       "对话成员必须大于0个",
	// ChangeSuperDto
	"ChangeSuperDto.Uid.required":   "请选择成员",
	"ChangeSuperDto.Super.required": "请选择要取消还是设置超级用户",
	"ChangeSuperDto.Super.min":      "超级用户数值应该是1或2",
	"ChangeSuperDto.Super.max":      "超级用户数值应该是1或2",

	// 工作流
	"WorkflowNodeQueryDto.TypeId.required":                "缺少工作流类型查询项",
	"WorkflowNodeDto.TypeId.required":                     "请选择工作流类型",
	"WorkflowNodeDto.Name.required":                       "请填写工作流节点名称",
	"WorkflowNodeDto.Node.min":                            "节点序号应该大于0",
	"WorkflowNodeDto.Node.max":                            "节点序号最大只允许999",
	"WorkflowNodeDto.End.oneof":                           "结束节点标记的值应该是0或1",
	"WorkflowNodeSaveCirculationDto.TypeId.required":      "缺少工作流类型",
	"WorkflowNodeSaveCirculationDto.Circulation.required": "缺少节点流转配置",
	"WorkflowNodeSaveCirculationDto.Circulation.gt":       "节点流转配置不能为空",

	// 其它
	"SingleUintRequired.ID.required": "缺少ID参数",
}

// FailHandle 表单验证失败字符串处理
func FailHandle(err error) string {
	// validator.ValidationErrors 是一个切片，不是真正的error，
	// 所以errors.Is和errors.As无法判断，只能用断言
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			// 只要第一个错误信息就行
			// example: []string{"WorkflowTypeDto.CirculationMode", "required"}
			return GetMapValue(strings.Join([]string{fieldError.Namespace(), fieldError.Tag()}, "."))
		}
	}

	var validatorInvalidValidationError *validator.InvalidValidationError
	if errors.As(err, &validatorInvalidValidationError) {
		logrus.Warnln("validator InvalidValidationError", err)
		return "参数校验失败，输入参数错误"
	}

	var jsonUnmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &jsonUnmarshalTypeError) {
		logrus.Errorln("json.UnmarshalTypeError", err)
		return "数据解析出错，请检查数据类型是否正确"
	}

	// 如果不符合以上的错误类型，就当作普通字符串处理
	// 不过上面校验了 validator.ValidationErrors 后，这里应该也不会去执行它了
	// 姑且认为是废弃吧，但是先保留 [Mr.Huan 2024-12-24 01:20]
	s := err.Error()
	// 记录日志
	logrus.Warnln("validator fail", s)
	// 生成正则对象
	reg := regexp.MustCompile("'([A-Za-z\\d.]+)'") // 取所有单引号内的内容
	res := reg.FindAllStringSubmatch(s, 3)         // 只匹配3次
	// [['UserRegisterForm.Password' UserRegisterForm.Password] ['Password' Password] ['required' required]]
	if len(res) < 3 {
		return "传递的参数格式有误" // 如果匹配结果小于3次，则代表匹配失败
	}
	textSlice := make([]string, 0)
	for i := 0; i < 3; i++ {
		if len(res[i]) < 2 {
			return "参数校验失败，未知格式: " + s // 表示返回的格式不是这里需要的
		}

		if i == 1 {
			continue // 第2次循环直接跳过，不需要这个值
		}

		// res[i] 第二个元素才是这里需要的
		textSlice = append(textSlice, res[i][1])
	}
	return GetMapValue(strings.Join(textSlice, "."))
}

// GetMapValue 获取表单校验Map的值
func GetMapValue(key string) string {
	if value, ok := formVerificationFailed[key]; ok {
		return value
	} else {
		return key
	}
}
