package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type redisconnectpool struct {
	pool          *redis.Pool
	redisip       string
	redispassword string
	redisname     string
}

func (redispool *redisconnectpool) ConnectRedis() {
	if nil == redispool.pool {
		redispool.pool = &redis.Pool{ //实例化一个连接池
			MaxIdle: 16, //最初的连接数量
			// MaxActive:1000000,    //最大连接数量
			MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
			IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
			Dial: func() (redis.Conn, error) { //要连接的redis数据库
				return redis.Dial("tcp", redispool.redisip, redis.DialPassword(redispool.redispassword))
			},
		}
	} else {
		fmt.Println("pool is not nil")
	}
}
func (redispool *redisconnectpool) SccredisSet(key string, value string) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	_, err := c.Do("Set", key, value)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (redispool *redisconnectpool) SccredisGet(key string) (value string, err error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	r, err := redis.String(c.Do("Get", key))
	if err != nil {
		fmt.Println("get key faild :", err)
		return "", err
	}
	return r, err
}
func (redispool *redisconnectpool) SccredisHSet(key string, filed string, value string) (err1 error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池
	_, err := c.Do("HSet", key, filed, value)
	//_, err = c.Do("HSet", "user01", "name1", "tom2")
	if err != nil {
		fmt.Println("hset key faild :", err)
		return err
	}
	return err
}
func (redispool *redisconnectpool) SccredisGetAll(key string) (value map[string]string, err1 error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	r, err := redis.StringMap(c.Do("HGETALL", key))
	if err != nil {
		fmt.Println("get key faild :", err)
		return nil, err
	}
	return r, err
}
func (redispool *redisconnectpool) SccredisHdel(key string, field string) (err error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	_, err = c.Do("HDEL", key, field)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (redispool *redisconnectpool) SccredisDel(key string) (err error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	_, err = c.Do("del", key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func sjdtst(s int) (out1 int) {
	out := 33
	return out
}
func main() {
	sjdtst(2)
}
