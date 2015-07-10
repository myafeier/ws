package lib

import (
	// "encoding/json"
	// "errors"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

type Redisc struct {
	pool       *redis.Pool
	dbNum      int
	connstring string
}

func Newredisc(conn string, num int) *Redisc {
	return &Redisc{
		dbNum:      num,
		connstring: conn,
	}
}

func (rs *Redisc) StartAndGc() error {
	rs.connectInit()
	c := rs.pool.Get()
	defer c.Close()
	return c.Err()
}

func (rs *Redisc) connectInit() {

	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rs.connstring)
		if err != nil {
			log.Fatalln("无法连接redis")
		}
		_, selecter := c.Do("select", rs.dbNum)
		if selecter != nil {
			c.Close()
			return nil, selecter
		}
		return
	}
	rs.pool = &redis.Pool{
		MaxIdle:     1,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func (rs *Redisc) Get(key string) (v interface{}, err error) {
	c := rs.pool.Get()
	defer c.Close()
	return c.Do("GET", key)
}
func (rs *Redisc) Pop(key string) (v interface{}, err error) {
	c := rs.pool.Get()
	defer c.Close()
	return c.Do("RPOP", key)
}

func (rs *Redisc) Do(commandName string, args ...interface{}) (replay interface{}, err error) {
	c := rs.pool.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}
