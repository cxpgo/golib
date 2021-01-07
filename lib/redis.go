package lib

import (
	"encoding/json"
	"errors"
	"reflect"
	"github.com/gomodule/redigo/redis"
	"fmt"
	"github.com/cxpgo/golib/model"
	"time"
)

// Redis client.
type Redis struct {
	pool *redis.Pool // Underlying connection pool.
	//group  string      // Configuration group.
	config *model.RedisConf // Configuration.
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

func InitRedis(config *model.RedisConf) {
	gRedis = New(config)
	_, err := gRedis.Do("Ping")
	if err == nil {
		Log.Info("===>Redis Init Successful<===")
	}
}

func New(config *model.RedisConf) *Redis {
	//fmt.Printf("redis_conf=%+v\n",config)
	// The MaxIdle is the most important attribute of the connection pool.
	// Only if this attribute is set, the created connections from client
	// can not exceed the limit of the server.
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
	//if r.group != "" {
	//	// If it is an instance object,
	//	// it needs to remove it from the instance Map.
	//	instances.Remove(r.group)
	//}
	//pools.Remove(fmt.Sprintf("%v", r.config))
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

// DoWithTimeout sends a command to the server and returns the received reply.
// The timeout overrides the read timeout set when dialing the connection.
func (c *Conn) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.do(timeout, commandName, args...)
}

/*func RedisConnFactory(name string) (redis.Conn, error) {
	if ConfRedisMap != nil && ConfRedisMap.List != nil {
		for confName, cfg := range ConfRedisMap.List {
			if name == confName {
				randHost := cfg.ProxyList[rand.Intn(len(cfg.ProxyList))]
				if cfg.ConnTimeout == 0 {
					cfg.ConnTimeout = 50
				}
				if cfg.ReadTimeout == 0 {
					cfg.ReadTimeout = 100
				}
				if cfg.WriteTimeout == 0 {
					cfg.WriteTimeout = 100
				}
				c, err := redis.Dial(
					"tcp",
					randHost,
					redis.DialConnectTimeout(time.Duration(cfg.ConnTimeout)*time.Millisecond),
					redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Millisecond))
				if err != nil {
					return nil, err
				}
				if cfg.Password != "" {
					if _, err := c.Do("AUTH", cfg.Password); err != nil {
						c.Close()
						return nil, err
					}
				}
				if cfg.Db != 0 {
					if _, err := c.Do("SELECT", cfg.Db); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, nil
			}
		}
	}
	return nil, errors.New("create redis conn fail")
}

func RedisLogDo(trace *TraceContext, c redis.Conn, commandName string, args ...interface{}) (interface{}, error) {
	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)

	endExecTime := time.Now()
	if err != nil {
		trace.LogTag = "_com_redis_failure"
		m := map[string]interface{}{
			"method":    commandName,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		}
		Log.Errorw(m,trace)
		//Log.TagError(trace, "_com_redis_failure", map[string]interface{}{
		//	"method":    commandName,
		//	"err":       err,
		//	"bind":      args,
		//	"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		//})
	} else {
		replyStr, _ := redis.String(reply, nil)
		trace.LogTag = "_com_redis_success"
		m := map[string]interface{}{
			"method":    commandName,
			"bind":      args,
			"reply":     replyStr,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		}
		Log.Infow(m,trace)
		//Log.Info(replyStr)
		//Log.TagInfo(trace, "_com_redis_success", map[string]interface{}{
		//	"method":    commandName,
		//	"bind":      args,
		//	"reply":     replyStr,
		//	"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		//})
	}
	return reply, err
}
//通过配置 执行redis
func RedisConfDo(trace *TraceContext, name string, commandName string, args ...interface{}) (interface{}, error) {
	c, err := RedisConnFactory(name)
	if err != nil {
		trace.LogTag = "_com_redis_failure"
		m := map[string]interface{}{
			"method": commandName,
			"err":    errors.New("RedisConnFactory_error:" + name),
			"bind":   args,
		}
		Log.Errorw(m,trace)

		//Log.TagError(trace, "_com_redis_failure", map[string]interface{}{
		//	"method": commandName,
		//	"err":    errors.New("RedisConnFactory_error:" + name),
		//	"bind":   args,
		//})
		return nil, err
	}
	defer c.Close()

	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		trace.LogTag = "_com_redis_failure"
		m := map[string]interface{}{
				"method":    commandName,
				"err":       err,
				"bind":      args,
				"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		}
		Log.Errorw(m,trace)
	} else {
		replyStr, _ := redis.String(reply, nil)
		trace.LogTag = "_com_redis_success"
		m := map[string]interface{}{
				"method":    commandName,
				"bind":      args,
				"reply":     replyStr,
				"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		}
		Log.Infow(m,trace)
		//Log.TagInfo(trace, "_com_redis_success", map[string]interface{}{
		//	"method":    commandName,
		//	"bind":      args,
		//	"reply":     replyStr,
		//	"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		//})
	}
	return reply, err
}
*/
