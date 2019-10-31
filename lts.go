package main

import (
	"conredis"
	"database/sql"
	"fmt"
	"log"
	"myconnsql"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/websocket"
)

var conn myconnsql.MyConn = new(myconnsql.Myconnection)
var sonn conredis.ConnRedis = new(conredis.MyConn)
var redisConn redis.Conn
var db *sql.DB
var i, j, k int = 0, 0, 0 //记录正在连接的个数 k用来记录登陆后退出的个数
var deluser map[int]int = make(map[int]int)
var users map[int]*websocket.Conn = make(map[int]*websocket.Conn) //用来接受连接的 用户
var usernames map[int]string = make(map[int]string)               //用来判断是否重复连接
var secords map[int]string = make(map[int]string)                 //用来存放读取的消息记录
var secordslen int                                                //用来记录数据个数
/*
	用来接受前端页面传来的用户信息用来注册用户信息
*/
func register(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pwd := r.FormValue("pwd")
	istrue := conn.Conn_select(name, pwd, db)
	// istrue := true
	if istrue && conn.Conn_insert(name, pwd, db) {
		// fmt.Println("注册成功")

		http.Redirect(w, r, "http://192.168.56.1:8989/LoginDemo/skip1.html", http.StatusFound)
	} else {
		http.Redirect(w, r, "http://192.168.56.1:8989/LoginDemo/skip2.html", http.StatusFound)
	}
}

/*
	接受用户信息完成登陆功能
*/
func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	pwd := r.FormValue("pwd")
	istrue, str := conn.Conn_sql(name, pwd, db)
	var k int
	for k = 0; k < j; k++ {
		if name == usernames[k] {
			istrue = false
			str = "该用户已经登陆"
		}
	}
	if istrue {
		usernames[j] = name
		// j = j + 1
		http.Redirect(w, r, "http://192.168.56.1:8989/LoginDemo/test.html", http.StatusFound)
	} else {
		fmt.Fprintf(w, str)
	}
}

/*
	服务器用来接收和传递信息
*/
func Echo(ws *websocket.Conn) {
	var err error
	users[i] = ws
	i = i + 1
	j = j + 1
	secords = sonn.Conn_List_Lrange("mysecord", redisConn)
	secordslen = len(secords)
	for {
		var j int
		var m int //用于循环
		var reply string
		var slen string
		var message string
		//for i := range c {
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			for j = 0; j < i; j++ {
				if users[j] == ws {
					if usernames[j] != "" {
						fmt.Println(usernames[j], "断开连接")
					}
					slen = strconv.Itoa(secordslen)
					// fmt.Println(slen)
					conn.Conn_update(usernames[j], slen, db)
					usernames[j] = ""
					deluser[k] = j
					k = k + 1
				}
			}
			break
		} else {
			var isuser bool = false //判断是否是合法用户
			for j = 0; j < i; j++ {
				if users[j] == ws && usernames[j] != "" {
					isuser = true
				}
			}
			if isuser {
				if reply == "#713" {
					reply = ""
					var secordslen int
					secords = sonn.Conn_List_Lrange("mysecord", redisConn)
					for secordslen = len(secords) - 1; secordslen >= 0; secordslen-- {
						reply += secords[secordslen]
					}
					for j = 0; j < i; j++ {
						if users[j] == ws {
							if err = websocket.Message.Send(users[j], reply); err != nil {
								fmt.Println("Can't send")
								break
							} else {
								fmt.Println(usernames[j], "查看了聊天记录")
								reply = ""
							}
						}
					}
				} else if reply == "#925" {
					// fmt.Println(reply)
					for j = 0; j < i; j++ {
						if users[j] == ws {
							reply = ""
							var secordslen int
							var slen int
							var namelen string
							secords = sonn.Conn_List_Lrange("mysecord", redisConn)
							namelen = conn.Conn_getlen(usernames[j], db)
							slen, err = strconv.Atoi(namelen)
							for secordslen = len(secords) - slen - 1; secordslen >= 0; secordslen-- {
								reply += secords[secordslen]
							}
							if err = websocket.Message.Send(users[j], reply); err != nil {
								fmt.Println("Can't send")
								break
							} else {
								fmt.Println(usernames[j], "查看了未读信息")
								reply = ""
							}
						}
					}
				} else {
					message = reply
					for j = 0; j < i; j++ {
						if users[j] == ws {
							reply = usernames[j]
						}
					}
					reply += "："
					reply += message
					fmt.Println(reply)
					reply += "&#13;&#10;"
					sonn.Conn_List_Lpush("mysecord", reply, redisConn)
					secordslen++
				}
			} else {
				reply = ""
			}
		}
		//	fmt.Println("Received back from client: " + reply)
		//}
		var msg string
		msg = reply      //接受信息
		var isfalse bool //用来标记该连接者是否退出
		// fmt.Println("Sending to client: " + msg)
		for j = 0; j < i; j++ {
			isfalse = true
			for m = 0; m < k; m++ {
				if j == deluser[m] {
					isfalse = false
				}
			}
			if isfalse {
				if err = websocket.Message.Send(users[j], msg); err != nil {
					fmt.Println("Can't send")
					break
				}
			}
		}
	}
}
func main() {
	var err error
	db, err = conn.Conn() //连接mysql
	if err != nil {
		log.Fatal("数据库连接失败")
	}
	redisConn, err = sonn.Conn() //连接redis
	if err != nil {
		log.Fatal("redis连接失败")
	}
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.Handle("/web", websocket.Handler(Echo))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("http listen failed")
	}
	conn.Conn_close(db)
}
