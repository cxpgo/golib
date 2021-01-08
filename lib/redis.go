package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cxpgo/golib/model/config"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"time"
)


var RedisMapPool map[string]*Redis
// Redis client.
type Redis struct {
	pool *redis.Pool // Underlying connection pool.
	//group  string      // Configuration group.
	config *config.RedisConf // Configuration.
}

// Redis connection.
type Conn struct {
	redis.Conn
}

// Pool statistics.
type PoolStats struct {
	redis.PoolStats
}

const (
	REDIS_POOL_IDLE_TIMEOUT  = 10
	REDIS_POOL_CONN_TIMEOUT  = 10
	REDIS_POOL_MAX_IDLE      = 10
	REDIS_POOL_MAX_ACTIVE    = 100
	REDIS_POOL_MAX_LIFE_TIME = 30
)

var (
// Pool map.
//pools = gmap.NewStrAnyMap(true)
)

func InitRedis(configs map[string]*config.RedisConf) {
	RedisMapPool = map[string]*Redis{}
	for confName,conf := range configs{
		redis := NewRedis(conf)
		RedisMapPool[confName] = redis
	}
	SetDefaultRedis("default")
	_, err := GRedis.Do("Ping")
	if err == nil {
		Log.Info("===>Redis Init Successful<===")
	}
}

func NewRedis(config *config.RedisConf) *Redis {

	if config.MaxIdle == 0 {
		config.MaxIdle = REDIS_POOL_MAX_IDLE
	}
	// This value SHOULD NOT exceed the connection limit of redis server.
	if config.MaxActive == 0 {
		config.MaxActive = REDIS_POOL_MAX_ACTIVE
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = REDIS_POOL_IDLE_TIMEOUT
	}
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = REDIS_POOL_CONN_TIMEOUT
	}
	if config.MaxConnLifetime == 0 {
		config.MaxConnLifetime = REDIS_POOL_MAX_LIFE_TIME
	}

	return &Redis{
		config: config,
		pool: &redis.Pool{
			Wait:            true,
			IdleTimeout:     time.Duration(config.IdleTimeout) * time.Second,
			MaxActive:       config.MaxActive,
			MaxIdle:         config.MaxIdle,
			MaxConnLifetime: time.Duration(config.MaxConnLifetime) * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial(
					"tcp",
					fmt.Sprintf("%s:%d", config.Host, config.Port),
					redis.DialConnectTimeout(time.Duration(config.ConnectTimeout)*time.Second),
					//redis.DialDatabase(config.Db),
					//redis.DialPassword(config.Password),
					//redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
					//redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
					//redis.DialUseTLS(config.TLS),
					//redis.DialTLSSkipVerify(config.TLSSkipVerify),
				)
				if err != nil {
					return nil, err
				}
				// AUTH
				if len(config.Password) > 0 {
					if _, err := c.Do("AUTH", config.Password); err != nil {
						return nil, err
					}
				}
				// DB
				if _, err := c.Do("SELECT", config.Db); err != nil {
					return nil, err
				}
				return c, nil
			},
			// After the conn is taken from the connection pool, to test if the connection is available,
			// If error is returned then it closes the connection object and recreate a new connection.
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}

}

func GetRedisByName(name string) *Redis{
	return RedisMapPool[name]
}
func SetDefaultRedis(name string){
	GRedis = RedisMapPool[name]
}

func CloseRedis()  {
	for _,redisPool := range RedisMapPool{
		redisPool.Close()
	}
}
// Do sends a command to the server and returns the received reply.
// Do automatically get a connection from pool, and close it when the reply received.
// It does not really "close" the connection, but drops it back to the connection pool.
func (r *Redis) Do(commandName string, args ...interface{}) (interface{}, error) {
	conn := &Conn{r.pool.Get()}
	defer conn.Close()
	return conn.Do(commandName, args...)
}

// DoWithTimeout sends a command to the server and returns the received reply.
// The timeout overrides the read timeout set when dialing the connection.
func (r *Redis) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (interface{}, error) {
	conn := &Conn{r.pool.Get()}
	defer conn.Close()
	return conn.DoWithTimeout(timeout, commandName, args...)
}

func (r *Redis) Close() error {
	return r.pool.Close()
}

func (r *Redis) Conn() *Conn {
	return &Conn{r.pool.Get()}
}

func (r *Redis) GetConn() *Conn {
	return r.Conn()
}

// Stats returns pool's statistics.
func (r *Redis) Stats() *PoolStats {
	return &PoolStats{r.pool.Stats()}
}

func (c *Conn) do(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	for k, v := range args {
		reflectValue = reflect.ValueOf(v)
		reflectKind = reflectValue.Kind()
		if reflectKind == reflect.Ptr {
			reflectValue = reflectValue.Elem()
			reflectKind = reflectValue.Kind()
		}
		switch reflectKind {
		case
			reflect.Struct,
			reflect.Map,
			reflect.Slice,
			reflect.Array:
			// Ignore slice type of: []byte.
			if _, ok := v.([]byte); !ok {
				if args[k], err = json.Marshal(v); err != nil {
					return nil, err
				}
			}
		}
	}
	if timeout > 0 {
		conn, ok := c.Conn.(redis.ConnWithTimeout)
		if !ok {
			return nil, errors.New(`current connection does not support "ConnWithTimeout"`)
		}
		return conn.DoWithTimeout(timeout, commandName, args...)
	}
	return c.Conn.Do(commandName, args...)
}

func (c *Conn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.do(0, commandName, args...)
}

func (c *Conn) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.do(timeout, commandName, args...)
}


