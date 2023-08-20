package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type Comment_respond struct {
	Com     Com `json:"comment"` // 评论成功返回评论内容，不需要重新拉取整个列表
	Respond Respond
}
type CommentList_respond struct {
	CommentList []Com `json:"comment_list"` // 评论列表
	Respond     Respond
}

// Comment
type Com struct {
	Content    string       `json:"content"`     // 评论内容
	CreateDate string       `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64        `json:"id"`          // 评论id
	Aut        User_respond `json:"user"`        // 评论用户信息
}

// 评论操作
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	actionType := c.Query("action_type")
	if len(token) == 0 {
		c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "账号不存在"})
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
			//获取用户信息
			var user User
			if err := tx.Preload("Videos").Take(&user, tk.ID).Error; err != nil {
				c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "账号不存在"})
				return err
			}
			var video Video
			if err := tx.Find(&video, video_id).Error; err != nil {
				c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "视频不存在"})
				return err
			}
			if actionType == "1" {
				comment_text := c.Query("comment_text")
				com := Comment{
					Con:        comment_text,
					CreateTime: time.Now().Format("01-02"),
					UID:        user.ID,
					VideoID:    video.ID,
				}
				if err := tx.Create(&com).Error; err != nil {
					c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "评论失败"})
					return err
				}
				c.JSON(http.StatusOK, Comment_respond{
					Com: Com{
						Content:    comment_text,
						CreateDate: com.CreateTime, // 评论发布日期，格式 mm-dd
						ID:         com.ID,         // 评论id
						Aut: User_respond{
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
						},
					},
					Respond: Respond{StatusCode: 0, StatusMsg: "评论成功"},
				})
			} else {
				comment_id := c.Query("comment_id")
				var com Comment
				if err := tx.Take(&com, comment_id).Error; err != nil {
					c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "删除评论失败"})
					return err
				}
				if err := tx.Delete(&com).Error; err != nil {
					c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "删除评论失败"})
					return err
				}
				c.JSON(http.StatusOK, Respond{StatusCode: 0, StatusMsg: "删除评论成功"})
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

// 查看视频的所有评论，按发布时间倒序
func CommentList(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	if len(token) == 0 {
		c.JSON(http.StatusOK, Respond{StatusCode: 0, StatusMsg: "请先登录"})
	} else {
		var tk Token
		if err := DB.Where("Token = ?", token).Take(&tk).Error; err != nil {
			c.JSON(http.StatusOK, Respond{
				StatusCode: 1,
				StatusMsg:  "请先登录",
			})
			return
		}
		var video Video
		DB.Preload("Comments").Take(&video, video_id)
		l := len(video.Comments)
		com := make([]Com, l)
		for i, x := range video.Comments {
			//获取该评论用户信息
			var user User
			DB.Preload("Videos").Take(&user, x.UID)
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
			com[l-i-1] = Com{
				Content:    x.Con,
				CreateDate: x.CreateTime,
				ID:         x.ID,
				Aut:        Aut,
			}
		}
		c.JSON(http.StatusOK, CommentList_respond{
			Respond:     Respond{StatusCode: 0},
			CommentList: com,
		})
	}
}
