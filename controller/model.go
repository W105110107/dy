package controller

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var DB *gorm.DB
var TimeLayoutStr = "2006-01-02 15:04:05"
var URL = "http://192.168.0.105:8080/static/"
var DUP = URL + "默认用户头像.jpg"      //默认用户头像
var DUHI = URL + "默认用户主页顶部大图.jpg" //默认用户主页顶部大图

type Iry []int64

// 视频
type Video struct {
	CoverURL     string    `gorm:"not null"` // 视频封面地址
	ID           int64     `gorm:"not null"` // 视频唯一标识
	PlayURL      string    `gorm:"not null"` // 视频播放地址
	Title        string    `gorm:"not null"` // 视频标题
	CreateTime   string    `gorm:"not null"` //创建时间
	Comments     []Comment //评论
	FavoriteList Iry       `gorm:"type:string"` //点赞用户列表
	UserID       int64     `gorm:"not null"`    //视频作者
}

// 用户
type User struct {
	Avatar          string //个人头像
	BackgroundImage string //个人主页顶部大图
	Signature       string // 个人简介
	TotalFavorite   int64  `gorm:"not null"`    // 获赞数量
	ID              int64  `gorm:"not null"`    // 用户id
	Name            string `gorm:"not null"`    // 用户名称
	FollowList      Iry    `gorm:"type:string"` //关注列表
	FollowerList    Iry    `gorm:"type:string"` //粉丝列表
	FavoriteList    Iry    `gorm:"type:string"` //点赞视频列表
	Videos          []Video
	PeopleMessages  []PeopleMessage
}

// 评论
type Comment struct {
	Con        string `gorm:"not null"` // 内容
	CreateTime string `gorm:"not null"` // 评论发布日期，格式 mm-dd
	ID         int64  `gorm:"not null"` // 评论id
	UID        int64  `gorm:"not null"` //评论用户id
	VideoID    int64  `gorm:"not null"` //评论所属视频
}

// 用户消息
type PeopleMessage struct {
	Con        string `gorm:"not null"` //消息内容
	ID         int64  `gorm:"not null"` // 消息ID
	UserID     int64  `gorm:"not null"` //发送消息的用户ID
	TUID       int64  `gorm:"not null"` //接受消息的用户ID
	CreateTime int64  `gorm:"not null"` // 消息发送时间时间戳
}

type Token struct {
	Token string `gorm:"not null"` // "用户名-密码"
	ID    int64  `gorm:"not null"` //对应用户id
}

type Respond struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// 自定义排序（按照时间递增）
type PM []PeopleMessage

func (M PM) Len() int {
	return len(M)
}
func (M PM) Swap(i, j int) {
	M[i], M[j] = M[j], M[i]
}
func (M PM) Less(i, j int) bool {
	return M[i].CreateTime < M[j].CreateTime
}

// 返回的视频流结构体
type Video_respond struct {
	Author        User_respond `json:"author"`         // 视频作者信息
	CommentCount  int64        `json:"comment_count"`  // 视频的评论总数
	CoverURL      string       `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64        `json:"favorite_count"` // 视频的点赞总数
	ID            int64        `json:"id"`             // 视频唯一标识
	IsFavorite    bool         `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string       `json:"play_url"`       // 视频播放地址
	Title         string       `json:"title"`          // 视频标题
}

// 返回的视频作者信息结构体
type User_respond struct {
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	ID              int64  `json:"id"`               // 用户id
	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Name            string `json:"name"`             // 用户名称
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  string `json:"total_favorited"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`       // 作品数
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

// Scan 从数据库中读取出来
func (i *Iry) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, i)
	return err
}

// Value 存入数据库
func (i Iry) Value() (driver.Value, error) {
	return json.Marshal(i)
}
