package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// 返回的结构体
type Feed_respond struct {
	NextTime   int64           `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64           `json:"status_code"` // 状态码，0-成功，其他值-失败
	VideoList  []Video_respond `json:"video_list"`  // 视频列表
}

// 是否关注或者是否点赞
func getBool(token *string, List *Iry) bool {
	if len(*token) == 0 {
		return false
	}
	var tk Token
	if err := DB.Where("Token = ?", token).Take(&tk).Error; err != nil {
		return false
	}
	tkid := tk.ID
	for _, x := range *List {
		if x == tkid {
			return true
		}
	}
	return false
}

func getTime(s *string) int64 {
	//获取传回客户端的时间戳
	loc, _ := time.LoadLocation("Local")                         //重要：获取时区
	timeObj, err := time.ParseInLocation(TimeLayoutStr, *s, loc) //指定日期 转 当地 日期对象 类型为 time.Time
	if err != nil {
		fmt.Println("parse time failed err :", err)
		return time.Now().Unix()
	}
	return timeObj.Unix()
}

// 视频流接口
func Feed(c *gin.Context) {
	input_time := c.Query("latest_time") //string类型
	var last_time time.Time
	if len(input_time) == 0 {
		// 处理传入的时间戳（这里是秒）
		temp, _ := strconv.ParseInt(input_time, 10, 64)
		last_time = time.Unix(temp, 0)
	} else {
		last_time = time.Now()
	}
	token := c.Query("token")
	var VideoList []Video
	DB.Where("CreateTime >= ?", last_time.Format(TimeLayoutStr)).Preload("Comments").Find(&VideoList)
	if len(VideoList) == 0 {
		DB.Preload("Comments").Find(&VideoList)
	}
	l := len(VideoList)
	Video_respondlist := make([]Video_respond, l)
	for i, x := range VideoList {
		//获取该视频作者信息
		var user User
		DB.Preload("Videos").Take(&user, x.UserID)
		Aut := User_respond{
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			FavoriteCount:   int64(len(user.FavoriteList)),             // 喜欢数
			FollowCount:     int64(len(user.FollowList)),               // 关注总数
			FollowerCount:   int64(len(user.FollowerList)),             // 粉丝总数
			ID:              user.ID,                                   // 用户id
			IsFollow:        getBool(&token, &(user.FollowerList)),     // true-已关注，false-未关注
			Name:            user.Name,                                 // 用户名称
			Signature:       user.Signature,                            // 个人简介
			TotalFavorited:  strconv.FormatInt(user.TotalFavorite, 10), // 获赞数量
			WorkCount:       int64(len(user.Videos)),                   // 作品数
		}
		Video_respondlist[l-i-1] = Video_respond{
			Author:        Aut,                                // 视频作者信息
			CommentCount:  int64(len(x.Comments)),             // 视频的评论总数
			CoverURL:      x.CoverURL,                         // 视频封面地址
			FavoriteCount: int64(len(x.FavoriteList)),         // 视频的点赞总数
			ID:            x.ID,                               // 视频唯一标识
			IsFavorite:    getBool(&token, &(x.FavoriteList)), // true-已点赞，false-未点赞
			PlayURL:       x.PlayURL,                          // 视频播放地址
			Title:         x.Title,                            // 视频标题
		}
	}

	c.JSON(http.StatusOK, Feed_respond{
		NextTime:   getTime(&VideoList[0].CreateTime),
		StatusCode: 0,
		VideoList:  Video_respondlist,
	})
}
