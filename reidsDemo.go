package main

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

type ConnRedis interface {
	/*
		用于连接redis，
		连接成功返回一个连接和nil
		连接失败返回一个nil和err
	*/
	Conn() (redis.Conn, error)
	/*
		key 表示储存的链表key
		message 表示存储的信息
	*/
	Conn_List_Lpush(key, message string, redisConn redis.Conn)
	/*
		key 表示需要连接的redis链表的名称
	*/
	Conn_List_Lrange(key string, redisConn redis.Conn) map[int]string
	Conn_List_LLen(key string, redisConn redis.Conn) int
	Conn_List_Lset(key string, index int, newvalues string, redisConn redis.Conn) bool
	/*
		关闭连接
	*/
	Conn_close(redisConn redis.Conn) //关闭redis
}
type MyConn struct {
}

/*
存储用户的聊天记录到redis指定列表中
*/
func (conn MyConn) Conn_List_Lpush(key, message string, redisConn redis.Conn) {
	_, err := redisConn.Do("lpush", key, message)
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
}

/*
根据参数key来查看列表中的数据
*/
func (conn MyConn) Conn_List_Lrange(key string, redisConn redis.Conn) map[int]string {
	values, err := redis.Values(redisConn.Do("lrange", key, "0", "1000"))
	if err != nil {
		fmt.Println("redis get failed:", err)
	}
	var eles map[int]string = make(map[int]string, 100000) //用于存储改key链表中的数据
	var i int = 0                                          // 用于key链表中数据的排序
	for _, v := range values {
		eles[i] = string(v.([]byte))
		i++
	}
	return eles
}
func (conn MyConn) Conn_List_LLen(key string, redisConn redis.Conn) int {
	values, err := redisConn.Do("llen", key)
	if err != nil {
		fmt.Println("redis llen failes", err)
	}
	index := values.(int)
	return index
}
func (conn MyConn) Conn_List_Lset(key string, index int, newvalues string, redisConn redis.Conn) bool {
	if index < 0 && index > conn.Conn_List_LLen(key, redisConn) {
		return false //自定义非0数为非法
	}
	_, err := redisConn.Do("lset", key, index, newvalues)
	if err != nil {
		fmt.Println("redis lset failes:", err)
	}
	return true
}

/*
用来连接redis，返回一个连接
*/
func (conn MyConn) Conn() (redis.Conn, error) {
	redisConn, err := redis.Dial("tcp", "192.168.56.1:6379")
	if err != nil {
		fmt.Println("连接失败:", err)
		return nil, err
	}

	return redisConn, nil
}

/*
关闭redis连接
*/
func (conn MyConn) Conn_close(redisConn redis.Conn) {
	defer redisConn.Close()
}

func main() {
	var redisConn redis.Conn
	var con ConnRedis = new(MyConn)
	redisConn, err := con.Conn()
	if err != nil {
		log.Fatal("连接失败")
	}
	con.Conn_List_Lset("mylion", 1, "lionada", redisConn)
	con.Conn_close(redisConn)
}
