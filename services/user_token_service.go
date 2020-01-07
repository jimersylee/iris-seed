package services

import (
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/datamodels"
	"github.com/jimersylee/iris-seed/repositories"
)

//UserTokenService处理用户数据模型的CRUID操作，
//它取决于用户存储库的操作。
//这是将数据源与更高级别的组件分离。
//因此，不同的存储库类型可以使用相同的逻辑，而无需任何更改。
//它是一个接口，它在任何地方都被用作接口
//因为我们可能需要在将来更改或尝试实验性的不同域逻辑。
type UserTokenServiceInterface interface {
	GetByID(id int64) *datamodels.UserToken
	DeleteByID(id int64) bool
}

var UserTokenService = NewUserTokenService(repositories.NewUserTokenRepository())

// NewUserTokenService返回默认用户服务
func NewUserTokenService(repo repositories.UserTokenRepository) *UserTokenServiceImpl {
	return &UserTokenServiceImpl{
		repo: repo,
	}
}

type UserTokenServiceImpl struct {
	repo repositories.UserTokenRepository
}

// DeleteByID按其id删除用户。
//如果删除则返回true，否则返回false。
func (s *UserTokenServiceImpl) DeleteByID(id int64) bool {
	s.repo.Delete(db.GetDB(), id)
	return true
}

func (s *UserTokenServiceImpl) GetByID(id int64) *datamodels.UserToken {
	return s.repo.FindOne(db.GetDB(), id)
}
