package main

import (
	m "dy/controller"
	"dy/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 连接数据库
func In() {
	username := "root"   //账号
	password := "123456" //密码
	host := "127.0.0.1"  //数据库地址，可以是Ip或者域名
	port := 3306         //数据库端口
	Dbname := "dy"       //数据库名
	timeout := "10s"     //连接超时，10秒

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）
		//如果没有这方面的要求，您可以在初始化时禁用它，这样可以获得60%的性能提升
		//SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "f_",  // 表名前缀
			SingularTable: true, // 单数表名
			NoLowerCase:   true, // 关闭小写转换
		},
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	// 连接成功
	m.DB = db
}

func main() {
	In()
	// 添加测试数据 (第一次启动是添加即可)
	//Add_Test_Data()

	go service.RunMessageServer()

	r := gin.Default()

	InitRouter(r)

	err := r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
}
