# 简易抖音

### 版本
go 1.17<br>
MySQL 8.0
### 获取第三方包
    go get -u github.com/u2takey/ffmpeg-go
    go get -u github.com/disintegration/imaging
    go get -u github.com/gin-gonic/gin
    go get gorm.io/driver/mysql
    go get gorm.io/gorm
### 改IP
![img.png](img.png)
将dy/controller/model 里红色框的ip改为自己电脑的ip地址
### 改数据库
在main.go函数里更改自己的数据库连接设置
### 启动
直接运行 main.go 即可
### 最后更新时间
2023-8-20