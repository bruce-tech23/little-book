package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"geektime-basic-learning2/little-book/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.key(uid)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (c *RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (c *RedisUserCache) Set(ctx context.Context, du domain.User) error {
	key := c.key(du.Id)
	d, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, d, c.expiration).Err()
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15, // user 专用的缓存，过期时间基本固定的，所以不需要外面传入
	}
}

// NewUserCacheV0 要想做到松耦合，一定不要自己去初始化你需要的东西，让外面传进来。尽量面向接口编程。
func NewUserCacheV0(addr string) *RedisUserCache {
	cmd := redis.NewClient(&redis.Options{Addr: addr}) // 这里初始化不太好，因为 Options 可能不只需要Addr,可能还需要别的参数
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
