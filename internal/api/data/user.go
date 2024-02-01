package data

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type UserRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (r *UserRepo) CreateUser(data *repo.User) error {
	return r.tx.Create(&data).Error
}

func (r *UserRepo) SaveUser(data *repo.User) error {
	return r.tx.Save(&data).Error
}

func (r *UserRepo) UpdatesUser(id uint64, data *repo.User) error {
	return r.tx.Model(&repo.User{}).Where("id = ?", id).Updates(&data).Error
}

func (r *UserRepo) DeleteUser(id uint64) error {
	return r.tx.Delete(&repo.User{}, id).Error
}

func (r *UserRepo) GetUser(id uint64) (*repo.User, error) {
	var user *repo.User
	err := r.tx.First(&user, id).Error
	return user, err
}

// Exist 用户是否存在
// 不记录错误也不返回错误
func (r *UserRepo) Exist(id uint64) bool {
	// 有记录就说明查到了
	return r.tx.Select("id").Where("id = ?", id).First(&repo.User{}).Error == nil
}

// ExistByUsername 用户是否存在
// 不记录错误也不返回错误
func (r *UserRepo) ExistByUsername(username string) bool {
	// 有记录就说明查到了
	return r.tx.Select("id").Where("user_login = ?", username).First(&repo.User{}).Error == nil
}

func (r *UserRepo) QueryUsernameAndPass(username, pwd string) (*repo.User, error) {
	var user *repo.User
	err := r.tx.Where("user_login = ?", username).Where("user_pass = ?", pwd).First(&user).Error
	return user, err
}

func (r *UserRepo) QueryUsername(username string) (*repo.User, error) {
	var user *repo.User
	err := r.tx.Where("user_login = ?", username).First(&user).Error
	return user, err
}

func (r *UserRepo) PageListUser(query dto.MemberListsQuery) ([]repo.User, int64, error) {
	var (
		count   int64
		members []repo.User
	)

	tx := r.tx.Model(&repo.User{})
	/* 查询 Start */
	if query.Id > 0 {
		// 如果是ID查询就只允许一个条件
		tx = tx.Where("id = ? ", query.Id)
	} else {
		if query.Status > 0 {
			tx = tx.Where("user_status = ? ", query.Status)
		}
	}

	if query.Username != "" {
		tx = tx.Where("user_login LIKE ?", "%"+query.Username+"%")
	}
	if query.Nickname != "" {
		tx = tx.Where("user_nickname LIKE ?", "%"+query.Nickname+"%")
	}
	if query.Mobile != "" {
		tx = tx.Where("mobile LIKE ?", "%"+query.Mobile+"%")
	}
	if query.Email != "" {
		tx = tx.Where("user_email LIKE ?", "%"+query.Email+"%")
	}
	/* 查询 End */

	// 统计数量
	err := tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询分页记录
	err = tx.Scopes(db.Paginate(&query.Page, &query.PageSize)).Find(&members).Error

	return members, count, err
}

func (r *UserRepo) SimpleList(key string) []dto.SimpleMemberList {
	var simpleList []dto.SimpleMemberList
	tx := r.tx.Model(&repo.User{}).Where("user_status = ?", 1)
	if key != "" {
		tx = tx.Where("user_login LIKE ? OR user_nickname LIKE ?", "%"+key+"%", "%"+key+"%")
	}
	tx.Find(&simpleList)
	return simpleList
}

func (r *UserRepo) UpdateUserStatus(id uint64, status int) error {
	return r.tx.Model(&repo.User{}).Where("id = ?", id).Updates(map[string]interface{}{"user_status": status}).Error
}

func (r *UserRepo) UpdateUserPass(id uint64, pwd string) error {
	// 同时更新修改密码的时间
	return r.tx.Model(&repo.User{}).Where("id = ?", id).Updates(repo.User{UserPass: pwd, LastEditPass: time.Now().Unix()}).Error
}

func (r *UserRepo) UpdateUserSuper(id uint64, super int8) error {
	return r.tx.Model(&repo.User{}).Where("id = ?", id).Updates(map[string]interface{}{"super": super}).Error
}

func (r *UserRepo) GetAdministrators() ([]repo.User, error) {
	var l []repo.User
	err := r.tx.Where("user_status = ?", 1).Where("super = ?", 1).Find(&l).Error
	return l, err
}

func NewUserRepo(tx *gorm.DB, ctx *gin.Context) repo.UserRepo {
	return &UserRepo{
		tx:  tx,
		ctx: ctx,
	}
}
