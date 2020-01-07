package repositories

import (
	"github.com/jimersylee/iris-seed/datamodels"
	"github.com/jinzhu/gorm"
)

// Query表示访问者和操作查询。
type Query func(datamodels.User) bool

// UserRepository处理用户实体/模型的基本操作。
//它是一个可测试的接口，即内存用户存储库或 连接到sql数据库。
type UserRepository interface {
	Exec(query Query, action Query, limit int, mode int) (ok bool)
	Select(query Query) (user datamodels.User, found bool)
	SelectMany(query Query, limit int) (results []datamodels.User)
	InsertOrUpdate(user datamodels.User) (updatedUser datamodels.User, err error)
	Delete(query Query, limit int) (deleted bool)
	FindByToken(db *gorm.DB, username string) *datamodels.User
	FindOne(db *gorm.DB, id int64) (user *datamodels.User)
}

// NewUserRepository返回一个新的基于用户内存的存储库，
//我们示例中唯一的存储库类型。
func NewUserRepository() UserRepository {
	return &userRepository{}
}

//userMemoryRepository是一个“UserRepository”
//使用内存数据源（map）管理用户。
type userRepository struct {
}

func (r *userRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
	return false
}

func (r *userRepository) FindOne(db *gorm.DB, id int64) (user *datamodels.User) {
	ret := &datamodels.User{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}
func (r *userRepository) FindByToken(db *gorm.DB, username string) *datamodels.User {
	ret := &datamodels.User{}
	if err := db.First(ret, "name= ? ", username).Error; err != nil {
		return nil
	}
	return ret
}

//Select接收查询方法
//为内部的每个用户模型触发查找我们想象中的数据源
//当该函数返回true时，它会停止迭代。
//它实际上是一个简单但非常游泳的原型函数
//自从我第一次想到它以来，我一直在使用它，
//希望你会发现它也很有用。
func (r *userRepository) Select(query Query) (user datamodels.User, found bool) {
	return
}

// SelectMany与Select相同但返回一个或多个datamodels.User作为切片。
//如果limit <= 0则返回所有内容。
func (r *userRepository) SelectMany(query Query, limit int) (results []datamodels.User) {
	return
}

// InsertOrUpdate将用户添加或更新到（内存）存储。
//返回新用户，如果有则返回错误。
func (r *userRepository) InsertOrUpdate(user datamodels.User) (datamodels.User, error) {
	return user, nil
}
func (r *userRepository) Delete(query Query, limit int) bool {
	return false
}
