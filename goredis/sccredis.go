package sccredis

import (
	"fmt"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

type Redisconnectpool struct {
	pool          *redis.Pool
	Redisip       string
	Redispassword string
	Redisname     string
}

func (redispool *Redisconnectpool) ConnectRedis() {
	if nil == redispool.pool {
		redispool.pool = &redis.Pool{ //实例化一个连接池
			MaxIdle: 16, //最初的连接数量
			// MaxActive:1000000,    //最大连接数量
			MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
			IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
			Dial: func() (redis.Conn, error) { //要连接的redis数据库
				return redis.Dial("tcp", redispool.Redisip, redis.DialPassword(redispool.Redispassword))
			},
		}
	} else {
		fmt.Println("pool is not nil")
	}
}
func (redispool *Redisconnectpool) SccredisSet(key string, value string) error {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	_, err := c.Do("Set", key, value)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}
func (redispool *Redisconnectpool) SccredisGet(key string) (string, error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	r, err := redis.String(c.Do("Get", key))
	if err != nil {
		//	fmt.Println("get key faild :", err)
		return "", err
	}
	return r, err
}

func (redispool *Redisconnectpool) SccredisHSet(key string, filed string, value string) error {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池
	_, err := c.Do("HSet", key, filed, value)
	//_, err = c.Do("HSet", "user01", "name1", "tom2")
	if err != nil {
		//fmt.Println("hset key faild :", err)
		return err
	}
	return err
}
func (redispool *Redisconnectpool) SccredisGetAll(key string) (map[string]string, error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	r, err := redis.StringMap(c.Do("HGETALL", key))
	if err != nil {
		//fmt.Println("get key faild :", err)
		return nil, err
	}
	return r, err
}

func (redispool *Redisconnectpool) SccredisBGetAll(key string) ([][]byte, error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	r, err := redis.ByteSlices(c.Do("HGETALL", key))
	if err != nil {
		//fmt.Println("get key faild :", err)
		return nil, err
	}
	return r, err
}
func (redispool *Redisconnectpool) SccredisHget(key string, filed string) (string, error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	r, err := redis.String(c.Do("HGET", key, filed))
	if err != nil {
		//fmt.Println("get key faild :", err)
		return "", err
	}
	return r, err
}

func (redispool *Redisconnectpool) SccredisBHget(key string, filed string) ([]byte, error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	r, err := redis.Bytes(c.Do("HGET", key, filed))
	if err != nil {
		//fmt.Println("get key faild :", err)
		return nil, err
	}
	return r, err
}
func (redispool *Redisconnectpool) SccredisHdel(key string, field string) error {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	_, err := c.Do("HDEL", key, field)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (redispool *Redisconnectpool) SccredisDel(key string) error {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	_, err := c.Do("del", key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (redispool *Redisconnectpool) SccredisAddgps(longitude string, latitude string, member string) error {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	_, err := c.Do("geoadd", "sccgps", longitude, latitude, member)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func (redispool *Redisconnectpool) SccredisDelgps(member string) error {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	_, err := c.Do("ZREM", "sccgps", member)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (redispool *Redisconnectpool) SccredisGetmembernearby(longitude string, latitude string, distance int) ([]string, error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池

	r, err := redis.Strings(c.Do("georadius", "sccgps", longitude, latitude, strconv.Itoa(distance), "m"))
	if err != nil {
		fmt.Println("get key faild :", err)
		return nil, err
	}
	return r, err
}
func (redispool *Redisconnectpool) SccredisGetmembergpsinf(members []string) ([]*[2]float64, error) {
	c := redispool.pool.Get() //从连接池，取一个链接
	defer c.Close()           //函数运行结束 ，把连接放回连接池
	tmosting := []interface{}{}
	tmosting = append(tmosting, "sccgps")
	for _, v := range members {
		tmosting = append(tmosting, v)
	}
	r, err := redis.Positions(c.Do("geopos", tmosting...))
	if err != nil {
		fmt.Println("get key faild :", err)
		return r, err
	}
	return r, err
}

/*func main() {
	var sjdtest Redisconnectpool
	sjdtest.redisip = "192.168.1.124:6379"
	sjdtest.redispassword = "123456"
	sjdtest.ConnectRedis()
	r, _ := sjdtest.SccredisGetmembernearby("118.32155", "31.123", 1000)
	sjdtest.SccredisGetmembergpsinf(r)
}*/
