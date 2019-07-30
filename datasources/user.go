package datasource

import (
	"github.com/jimersylee/iris-seed/datamodels"
	"github.com/jinzhu/gorm"
)

//引擎来自何处获取数据，在这种情况下是用户。
type Engine uint32

const (
	//内存代表简单的内存位置;
	// map[int64]datamodels.User随时可以使用，这是我们在这个例子中的来源。
	Memory Engine = iota
	// Bolt for boltdb source location.
	Bolt
	// MySQL for mysql-compatible source location.
	MySQL
)

var UserRepository = newUserRepository()

func newUserRepository() *userRepository {
	return &userRepository{}
}

type userRepository struct {
}

//根据id找用户
func (this *userRepository) Get(db *gorm.DB, id int64) *datamodels.Users {
	ret := &datamodels.Users{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}
