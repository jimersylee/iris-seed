package datasource

import (
	"errors"
	"github.com/jimersylee/iris-seed/datamodels"
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

//为了简单起见，Load Users从内存中返回所有用户（空map）。
func LoadUsers(engine Engine) (map[int64]datamodels.User, error) {
	if engine != Memory {
		return nil, errors.New("for the shake of simplicity we're using a simple map as the data source")
	}
	return make(map[int64]datamodels.User), nil
}
