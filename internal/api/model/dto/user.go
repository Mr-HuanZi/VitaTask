package dto

type PostUid struct {
	Uid uint64 `json:"uid" binding:"required"`
}

type UserRegisterForm struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required,gte=5,lte=20,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
	UserNickname    string `json:"userNickname" binding:"required"`
	UserEmail       string `json:"userEmail" binding:"omitempty,email"` // omitempty 省略空值
	Mobile          string `json:"mobile"`
}

type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code"`
}

type UserInfoDto struct {
	UserNickname string `json:"nickname" binding:"required"`
	UserEmail    string `json:"email,omitempty" binding:"omitempty,email"`
	Mobile       string `json:"mobile,omitempty"`
	Signature    string `json:"signature,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	Sex          int8   `json:"sex,omitempty"`
	Birthday     string `json:"birthday,omitempty"`
}

type ChangePasswordDto struct {
	OldPassword     string `json:"old_password" binding:"required"`
	Password        string `json:"password" binding:"required,gte=5,lte=20,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type ChangeSuperDto struct {
	Uid   uint64 `json:"uid" binding:"required"`
	Super int8   `json:"super" binding:"required,min=1,max=2"`
}
