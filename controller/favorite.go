package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type FavoriteList_respond struct {
	Respond   Respond
	VideoList []Video_respond `json:"video_list"` // 用户点赞视频列表
}

// 赞操作
func FavoriteAction(c *gin.Context) {
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	token := c.Query("token")
	if len(token) == 0 {
		c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "请先登录在点赞"})
	} else {
		var tk Token
		if err := DB.Where("Token = ?", token).Take(&tk).Error; err != nil {
			c.JSON(http.StatusOK, Respond{
				StatusCode: 1,
				StatusMsg:  "请先登录",
			})
			return
		}
		err := DB.Transaction(func(tx *gorm.DB) error {
			//获取当前登录用户id
			var vid, _ = strconv.ParseInt(video_id, 10, 64)
			var video Video
			if err := tx.Take(&video, vid).Error; err != nil {
				c.JSON(http.StatusOK, Respond{
					StatusCode: 1,
					StatusMsg:  "点赞失败",
				})
				return err
			}
			if action_type == "1" {
				//在该视频的点赞列表里添加点赞用户的id
				video.FavoriteList = append(video.FavoriteList, tk.ID)
				if err := tx.Select("FavoriteList").Save(&video).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "点赞失败",
					})
					return err
				}
				//视频作者的获赞数量 加1
				var aut User
				if err := tx.Take(&aut, video.UserID).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "点赞失败",
					})
					return err
				}
				aut.TotalFavorite++
				if err := tx.Select("TotalFavorite").Save(&aut).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "点赞失败",
					})
					return err
				}
				//将该视频id添加到用户的点赞（喜欢）列表里
				var user User
				if err := tx.Take(&user, tk.ID).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "点赞失败",
					})
					return err
				}
				user.FavoriteList = append(user.FavoriteList, vid)
				if err := tx.Select("FavoriteList").Save(&user).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "点赞失败",
					})
					return err
				}
				c.JSON(http.StatusOK, Respond{StatusCode: 0, StatusMsg: "点赞成功"})
			} else {
				//在该视频的点赞列表里删除点赞用户的id
				index := -1
				for i, x := range video.FavoriteList {
					if x == tk.ID {
						index = i
						break
					}
				}
				if index == -1 {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "取消点赞失败",
					})
					return errors.New("取消点赞失败")
				}
				video.FavoriteList = append(video.FavoriteList[:index], video.FavoriteList[index+1:]...)
				if err := tx.Select("FavoriteList").Save(&video).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "取消点赞失败",
					})
					return err
				}
				//视频作者的获赞数量 减1
				var aut User
				if err := tx.Take(&aut, video.UserID).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "取消点赞失败",
					})
					return err
				}
				aut.TotalFavorite--
				if err := tx.Select("TotalFavorite").Save(&aut).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "取消点赞失败",
					})
					return err
				}
				//将该视频id从用户的点赞（喜欢）列表里删除
				var user User
				if err := tx.Take(&user, tk.ID).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "取消点赞失败",
					})
					return err
				}
				index = -1
				for i, x := range user.FavoriteList {
					if x == vid {
						index = i
						break
					}
				}
				if index == -1 {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "取消点赞失败",
					})
					return errors.New("取消点赞失败")
				}
				user.FavoriteList = append(user.FavoriteList[:index], user.FavoriteList[index+1:]...)
				if err := tx.Select("FavoriteList").Save(&user).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "取消点赞失败",
					})
					return err
				}
				c.JSON(http.StatusOK, Respond{StatusCode: 0, StatusMsg: "取消点赞成功"})
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

// 用户的所有点赞视频
func FavoriteList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	var user User
	DB.Take(&user, user_id)

	var VideoList []Video
	In := make([]int64, len(user.FavoriteList))
	for i, x := range user.FavoriteList {
		In[i] = x
	}
	DB.Find(&VideoList, "ID in (?)", In)

	l := len(VideoList)
	Video_respondlist := make([]Video_respond, l)

	for i, x := range VideoList {

		//获取该视频作者信息
		var a User
		DB.Where("ID = ?", x.UserID).Preload("Videos").Find(&a)
		Aut := User_respond{
			Avatar:          a.Avatar,
			BackgroundImage: a.BackgroundImage,
			ID:              a.ID,                                   // 用户id
			Name:            a.Name,                                 // 用户名称
			FollowCount:     int64(len(a.FollowList)),               // 关注总数
			FollowerCount:   int64(len(a.FollowerList)),             // 粉丝总数
			IsFollow:        getBool(&token, &(a.FollowerList)),     // true-已关注，false-未关注
			Signature:       a.Signature,                            // 个人简介
			TotalFavorited:  strconv.FormatInt(a.TotalFavorite, 10), // 获赞数量
			WorkCount:       int64(len(a.Videos)),                   // 作品数
			FavoriteCount:   int64(len(a.FavoriteList)),             // 喜欢数
		}
		Video_respondlist[l-i-1] = Video_respond{
			ID:            x.ID,                               // 视频唯一标识
			Author:        Aut,                                // 视频作者信息
			PlayURL:       x.PlayURL,                          // 视频播放地址
			CoverURL:      x.CoverURL,                         // 视频封面地址
			FavoriteCount: int64(len(x.FavoriteList)),         // 视频的点赞总数
			CommentCount:  int64(len(x.Comments)),             // 视频的评论总数
			IsFavorite:    getBool(&token, &(x.FavoriteList)), // true-已点赞，false-未点赞
			Title:         x.Title,                            // 视频标题
		}
	}
	c.JSON(http.StatusOK, FavoriteList_respond{
		Respond:   Respond{StatusCode: 0},
		VideoList: Video_respondlist,
	})
}
