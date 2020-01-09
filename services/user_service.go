package services

import (
	"errors"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/datamodels"
	"github.com/jimersylee/iris-seed/repositories"
)

//UserService处理用户数据模型的CRUID操作，
//它取决于用户存储库的操作。
//这是将数据源与更高级别的组件分离。
//因此，不同的存储库类型可以使用相同的逻辑，而无需任何更改。
//它是一个接口，它在任何地方都被用作接口
//因为我们可能需要在将来更改或尝试实验性的不同域逻辑。
type UserServiceInterface interface {
	GetAll() []datamodels.User
	GetByID(id int64) datamodels.User
	GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool)
	DeleteByID(id int64) bool
	Update(id int64, user datamodels.User) (datamodels.User, error)
	UpdatePassword(id int64, newPassword string) (datamodels.User, error)
	UpdateUsername(id int64, newUsername string) (datamodels.User, error)
	Create(userPassword string, user datamodels.User) (datamodels.User, error)
}

var UserService = NewUserService(repositories.NewUserRepository())

// NewUserService返回默认用户服务
func NewUserService(repo repositories.UserRepository) *userServiceImpl {
	return &userServiceImpl{
		repo: repo,
	}
}

type userServiceImpl struct {
	repo repositories.UserRepository
}

// GetAll返回所有用户。
func (s *userServiceImpl) GetAll() []datamodels.User {
	return s.repo.SelectMany(func(_ datamodels.User) bool {
		return true
	}, -1)
}

// GetByID根据其id返回用户。
func (s *userServiceImpl) GetByID(id int64) *datamodels.User {

	return s.repo.FindOne(db.GetDB(), id)
}

//获取yUsernameAndPassword根据用户名和密码返回用户，
//用于身份验证。
func (s *userServiceImpl) GetUserByUsernameAndPassword(username, userPassword string) (datamodels.User, bool) {
	if username == "" || userPassword == "" {
		return datamodels.User{}, false
	}
	user := s.repo.FindByUserName(db.GetDB(), username)
	if user != nil {
		if ok, _ := datamodels.ValidatePassword(userPassword, []byte(user.Password)); ok {
			return *user, true
		}
		return datamodels.User{}, false
	}
	return datamodels.User{}, false
}

//获取yUsernameAndPassword根据用户名和密码返回用户，
//用于身份验证。
func (s *userServiceImpl) GetByUsername(username string) (datamodels.User, bool) {
	user := s.repo.FindByUserName(db.GetDB(), username)
	if user != nil {
		return *user, true
	}
	return datamodels.User{}, false
}

//更新现有用户的每个字段的更新，
//通过公共API使用是不安全的
//但是我们将在web  controllers/user_controller.go#PutBy上使用它
//为了向您展示它是如何工作的。
func (s *userServiceImpl) Update(id int64, user datamodels.User) (err error) {
	user.ID = id
	return s.repo.InsertOrUpdate(db.GetDB(), &user)
}

// UpdatePassword更新用户的密码。
func (s *userServiceImpl) UpdatePassword(id int64, newPassword string) (datamodels.User, error) {
	////更新用户并将其返回。
	//hashed, err := datamodels.GeneratePassword(newPassword)
	//if err != nil {
	//	return datamodels.User{}, err
	//}
	//return s.Update(id, datamodels.User{
	//	HashedPassword: hashed,
	//})
	return datamodels.User{}, nil
}

// UpdateUsername更新用户的用户名
func (s *userServiceImpl) UpdateUsername(id int64, newUsername string) (err error) {
	return s.Update(id, datamodels.User{
		Name: newUsername,
	})
}

//创建插入新用户，
// userPassword是客户端类型的密码
//它将在插入我们的存储库之前进行哈希处理
func (s *userServiceImpl) Create(userPassword string, user datamodels.User) (err error) {
	if user.ID > 0 || userPassword == "" || user.Name == "" {
		return errors.New("unable to create this user")
	}
	hashed, err := datamodels.GeneratePassword(userPassword)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return s.repo.InsertOrUpdate(db.GetDB(), &user)
}

// DeleteByID按其id删除用户。
//如果删除则返回true，否则返回false。
func (s *userServiceImpl) DeleteByID(id int64) bool {
	return s.repo.Delete(func(m datamodels.User) bool {
		return m.ID == id
	}, 1)
}
