package repositories

import (
	"errors"
	"github.com/jimersylee/iris-seed/datamodels"
	"sync"
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
}

// NewUserRepository返回一个新的基于用户内存的存储库，
//我们示例中唯一的存储库类型。
func NewUserRepository(source map[int64]datamodels.User) UserRepository {
	return &userMemoryRepository{source: source}
}

//userMemoryRepository是一个“UserRepository”
//使用内存数据源（map）管理用户。
type userMemoryRepository struct {
	source map[int64]datamodels.User
	mu     sync.RWMutex
}

const (
	// ReadOnlyMode将RLock(read) 数据。
	ReadOnlyMode = iota
	// ReadWriteMode将锁定(read/write)数据。
	ReadWriteMode
)

func (r *userMemoryRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
	loops := 0
	if mode == ReadOnlyMode {
		r.mu.RLock()
		defer r.mu.RUnlock()
	} else {
		r.mu.Lock()
		defer r.mu.Unlock()
	}
	for _, user := range r.source {
		ok = query(user)
		if ok {
			if action(user) {
				loops++
				if actionLimit >= loops {
					break // break
				}
			}
		}
	}
	return
}

//Select接收查询方法
//为内部的每个用户模型触发查找我们想象中的数据源
//当该函数返回true时，它会停止迭代。
//它实际上是一个简单但非常游泳的原型函数
//自从我第一次想到它以来，我一直在使用它，
//希望你会发现它也很有用。
func (r *userMemoryRepository) Select(query Query) (user datamodels.User, found bool) {
	found = r.Exec(query, func(m datamodels.User) bool {
		user = m
		return true
	}, 1, ReadOnlyMode)
	//设置一个空的datamodels.User，如果根本找不到的话
	if !found {
		user = datamodels.User{}
	}
	return
}

// SelectMany与Select相同但返回一个或多个datamodels.User作为切片。
//如果limit <= 0则返回所有内容。
func (r *userMemoryRepository) SelectMany(query Query, limit int) (results []datamodels.User) {
	r.Exec(query, func(m datamodels.User) bool {
		results = append(results, m)
		return true
	}, limit, ReadOnlyMode)
	return
}

// InsertOrUpdate将用户添加或更新到（内存）存储。
//返回新用户，如果有则返回错误。
func (r *userMemoryRepository) InsertOrUpdate(user datamodels.User) (datamodels.User, error) {
	id := user.ID
	if id == 0 {
		var lastID int64
		//找到最大的ID，以便不重复 在制作应用中，
		//您可以使用第三方库以生成UUID作为字符串。
		r.mu.RLock()
		for _, item := range r.source {
			if item.ID > lastID {
				lastID = item.ID
			}
		}
		r.mu.RUnlock()
		id = lastID + 1
		user.ID = id
		r.mu.Lock()
		r.source[id] = user
		r.mu.Unlock()
		return user, nil
	}
	//基于user.ID更新操作，
	//这里我们将允许更新海报和流派，如果不是空的话。
	//或者我们可以做替换：
	// r.source [id] =user
	//的代码;
	current, exists := r.Select(func(m datamodels.User) bool {
		return m.ID == id
	})
	if !exists { // ID不是真实的，返回错误。
		return datamodels.User{}, errors.New("failed to update a nonexistent user")
	}
	//和r.source [id] =user 进行纯替换
	if user.Username != "" {
		current.Username = user.Username
	}
	if user.Firstname != "" {
		current.Firstname = user.Firstname
	}
	r.mu.Lock()
	r.source[id] = current
	r.mu.Unlock()
	return user, nil
}
func (r *userMemoryRepository) Delete(query Query, limit int) bool {
	return r.Exec(query, func(m datamodels.User) bool {
		delete(r.source, m.ID)
		return true
	}, limit, ReadWriteMode)
}
