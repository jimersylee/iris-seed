package repositories

import (
	"github.com/jimersylee/go-steam-proxy/models"
	"github.com/jinzhu/gorm"
)

// Query表示访问者和操作查询。
type Query func(models.User) bool

// UserRepository处理用户实体/模型的基本操作。
//它是一个可测试的接口，即内存用户存储库或 连接到sql数据库。
type UserRepository interface {
	Select(query Query) (user models.User, found bool)
	SelectMany(query Query, limit int) (results []models.User)
	InsertOrUpdate(db *gorm.DB, user *models.User) (err error)
	Delete(query Query, limit int) (deleted bool)
	FindByUserName(db *gorm.DB, username string) *models.User
	FindOne(db *gorm.DB, id int64) (user *models.User)
}

// NewUserRepository返回一个新的基于mysql的存储库，
//我们示例中唯一的存储库类型。
func NewUserRepository() UserRepository {
	return &userRepository{}
}

//userMemoryRepository是一个“UserRepository”
//使用内存数据源（map）管理用户。
type userRepository struct {
}

func (r *userRepository) FindOne(db *gorm.DB, id int64) (user *models.User) {
	ret := &models.User{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}
func (r *userRepository) FindByUserName(db *gorm.DB, username string) *models.User {
	ret := &models.User{}
	if err := db.First(ret, "name= ? ", username).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userRepository) Select(query Query) (user models.User, found bool) {
	return
}

// SelectMany与Select相同但返回一个或多个datamodels.User作为切片。
//如果limit <= 0则返回所有内容。
func (r *userRepository) SelectMany(query Query, limit int) (results []models.User) {
	return
}

// InsertOrUpdate将用户添加或更新到（内存）存储。
//返回新用户，如果有则返回错误。
func (r *userRepository) InsertOrUpdate(db *gorm.DB, user *models.User) (err error) {
	return db.Save(user).Error
}
func (r *userRepository) Delete(query Query, limit int) bool {
	return false
}
