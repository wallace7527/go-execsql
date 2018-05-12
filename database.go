/*
Copyright 2018 The go-execsql Authors
This file is part of the go-execsql library.

database
封装数据库相关操作


wanglei.ok@foxmail.com

1.0
版本时间：16:22 2018/5/9

*/

package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"bytes"
)

//数据库操作对象
var db *sql.DB

const (
	//连接池属性
	POOL_MAXOPENCONNS = 10	//最大连接数
	POOL_MAXIDLECONNS = 2	//空闲连接数
)

//打开数据库
func OpenDatabase(dsn string) error {
	db1, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	//连接池
	db1.SetMaxOpenConns(POOL_MAXOPENCONNS)
	db1.SetMaxIdleConns(POOL_MAXIDLECONNS)
	//连接
	if  err = db1.Ping(); err != nil {
		return err
	}
	db = db1
	return nil
}

//关闭数据库
func CloseDatabase() {
	db.Close()
}

//自定义事务结构
type MyTx struct {
	Tx *sql.Tx
}

//开始事务并返回事务对象
func TxBegin() (*MyTx, error) {
	tx, err := db.Begin()
	return &MyTx{tx}, err
}

//提交事务
func (x *MyTx)Commit() error {
	return x.Tx.Commit()
}

//回滚事务
func (x *MyTx)Rollback() error {
	return x.Tx.Rollback()
}


//执行一条SQL语句
func (x *MyTx) ExecSQL(sql string) (sql.Result, error) {
	return x.Tx.Exec(sql)
}

//执行一条SQL语句
func ExecSQL(sql string) (sql.Result, error) {
	return db.Exec(sql)
}

func ResultString(r sql.Result) string {
	var buf bytes.Buffer

	rowsAffected, err := r.RowsAffected()
	buf.WriteString("RowsAffected:")
	if err != nil {
		buf.WriteString(err.Error())
	}else {
		buf.WriteString(strconv.FormatInt(rowsAffected,10))
	}

	lastInsertId, err := r.LastInsertId()
	buf.WriteString(", LastInsertId:")
	if err != nil {
		buf.WriteString(err.Error())
	}else {
		buf.WriteString(strconv.FormatInt(lastInsertId,10))
	}
	return  buf.String()
}