package db

import "fmt"

var DB DataBases

type DataBases struct {
	MysqlDB mysqlDB // 数据库连接
}

func init() {
	fmt.Println("init mysql redis mongo ")
	// 初始化数据库
	initMysqlDB()
}
