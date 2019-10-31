package myconnsql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MyConn interface {
	Conn() (*sql.DB, error)                                       //用于数据库连接
	Conn_select(username, userpwd string, db *sql.DB) bool        //用户查询用户表中是否有该用户
	Conn_insert(username, userpwd string, db *sql.DB) bool        //用于注册用户信息插入数据
	Conn_sql(username, userpwd string, db *sql.DB) (bool, string) //用来判断用户名密码
	Conn_update(username string, slen string, db *sql.DB) bool    //更新namesecord表中数据，来操作未读信息
	Conn_getlen(username string, db *sql.DB) string               //得到用户已读信息的位置
	Conn_close(db *sql.DB)                                        //关闭数据库
}

type Myconnection struct {
}

var id int64

/*
用于注册信息判断该用户是否已经存在
*/
func (conn Myconnection) Conn_select(username, userpwd string, db *sql.DB) bool {
	sql := "select * from t_user"
	stmt, err := db.Prepare(sql)
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("数据库查询错误")
	}
	for rows.Next() {
		var id string
		var name string
		var pwd string
		rows.Scan(&id, &name, &pwd)
		if username == name {
			return false
		}
	}
	return true
}

/*
1、用于更新用户信息，插入新的用户信息
2、用来更新用户查看消息记录信息中用户信息，刚注册的用户默认从头开始查看未读信息
*/
func (conn Myconnection) Conn_insert(username, userpwd string, db *sql.DB) bool {
	sql := "insert into t_user(name,password) values(?,?)"
	stmt, err := db.Prepare(sql)
	res, err := stmt.Exec(username, userpwd)
	sql = "insert into namesecord values(?,'0')"
	stmt, err = db.Prepare(sql)
	res, err = stmt.Exec(username)
	ins_id, err := res.LastInsertId()
	id = ins_id
	if err != nil {
		log.Fatal("数据库插入错误连接失败")
	}
	return true
}

/*
用来判断该用户是否存在，并验证用户名密码
*/
func (conn Myconnection) Conn_sql(username, userpwd string, db *sql.DB) (bool, string) {
	sql := "select * from t_user"
	stmt, err := db.Prepare(sql)
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("数据库查询失败")
	}
	for rows.Next() {
		var id string
		var name string
		var pwd string
		rows.Scan(&id, &name, &pwd)
		if username == name && userpwd == pwd {
			return true, "登陆成功"
		}
	}
	return false, "用户名或密码错误"
}

/*
用来更新用户查看消息记录
*/
func (con Myconnection) Conn_update(username string, slen string, db *sql.DB) bool {
	sql := "update namesecord set slen = ? where name = ?"
	stmt, err := db.Prepare(sql)
	res, err := stmt.Exec(slen, username)
	ins_id, err := res.LastInsertId()
	id = ins_id
	if err != nil {
		log.Fatal("数据库更新失败")
	}
	return true
}

/*
用来得到用户已经浏览的消息位置
*/
func (con Myconnection) Conn_getlen(username string, db *sql.DB) string {
	sql := "select slen from namesecord where name = ?"
	var slen string
	stmt, err := db.Prepare(sql)
	err = stmt.QueryRow(username).Scan(&slen)
	if err != nil {
		log.Fatal("namesecord数据库查询失败")
	}
	return slen
}

/*
用来连接数据库，并返回一个连接
*/
func (conn Myconnection) Conn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:1234@tcp(192.168.56.1:3306)/lion")
	if err != nil {
		log.Fatal("连接失败")
		return nil, err
	}
	return db, nil
}

/*
用来关闭数据库连接
*/
func (conn Myconnection) Conn_close(db *sql.DB) {
	db.Close()
}
