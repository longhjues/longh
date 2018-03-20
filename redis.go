package longh

import (
	"encoding/json"
	"strings"

	"github.com/garyburd/redigo/redis"
)

// NewRedisPool 新建一个Redis连接池 URL优先
func NewRedisPool(addr, passwd string, db int) *redis.Pool {
	b := strings.HasPrefix(addr, "redis://")
	var dialFunc func() (redis.Conn, error)
	switch {
	case b && passwd == "":
		dialFunc = func() (redis.Conn, error) {
			return redis.DialURL(addr, redis.DialDatabase(db))
		}
	case b && passwd != "":
		dialFunc = func() (redis.Conn, error) {
			return redis.DialURL(addr, redis.DialDatabase(db), redis.DialPassword(passwd))
		}
	case !b && passwd == "":
		dialFunc = func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, redis.DialDatabase(db))
		}
	case !b && passwd != "":
		dialFunc = func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, redis.DialDatabase(db), redis.DialPassword(passwd))
		}
	}

	return &redis.Pool{
		MaxIdle:   10,
		MaxActive: 200,
		Dial: func() (redis.Conn, error) {
			c, err := dialFunc()
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// BLPOPUnmarshalJSON 从队列阻塞获取json数据并解析
func BLPOPUnmarshalJSON(conn redis.Conn, key string, data interface{}) error {
	b, err := redis.Bytes(conn.Do("BLPOP", key, 0))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, data)
}

// LPUSHMarshalJSON json数据编码并存入队列
func LPUSHMarshalJSON(conn redis.Conn, key string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("LPUSH", key, b)
	return err
}
