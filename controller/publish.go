package controller

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type VideoList_response struct {
	Respond   Respond
	VideoList []Video_respond `json:"video_list"` // 用户发布的视频列表
}

// 第一个参数是 视频的相对地址，第二个参数是截出来的图片的名称
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (err error) {
	snapshotPath = "./public/" + snapshotPath
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	err = imaging.Save(img, snapshotPath+".jpg")
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}
	return nil
}

// 发布视频
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	if len(token) == 0 {
		c.JSON(http.StatusOK, Respond{
			StatusCode: 1,
			StatusMsg:  "请先登录",
		})
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
			data, err := c.FormFile("data")
			if err != nil {
				c.JSON(http.StatusOK, Respond{
					StatusCode: 1,
					StatusMsg:  "发布失败",
				})
				return err
			}
			var user User
			if err = tx.Preload("Videos").Take(&user, tk.ID).Error; err != nil {
				c.JSON(http.StatusOK, Respond{
					StatusCode: 1,
					StatusMsg:  "发布失败",
				})
				return err
			}
			l := len(user.Videos)
			finalName := fmt.Sprintf("%d_%d_%s", user.ID, l+1, title)
			err = c.SaveUploadedFile(data, "./public/"+finalName+".mp4")
			if err != nil {
				c.JSON(http.StatusOK, Respond{
					StatusCode: 1,
					StatusMsg:  "发布失败",
				})
				return err
			} else {
				err = GetSnapshot("./public/"+finalName+".mp4", finalName, 1)
				if err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "发布失败",
					})
					return err
				} else {
					video := Video{
						CoverURL:     URL + finalName + ".jpg",
						PlayURL:      URL + finalName + ".mp4",
						Title:        title,
						CreateTime:   time.Now().Format(TimeLayoutStr),
						Comments:     []Comment{},
						FavoriteList: Iry{},
						UserID:       user.ID,
					}
					if err = tx.Create(&video).Error; err != nil {
						c.JSON(http.StatusOK, Respond{
							StatusCode: 1,
							StatusMsg:  err.Error(),
						})
						return err
					}
					c.JSON(http.StatusOK, Respond{
						StatusCode: 0,
						StatusMsg:  "上传成功",
					})
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

// 用户的视频发布列表，直接列出用户所有投稿过的视频
func PublishList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")

	var user User
	err := DB.Preload("Videos").Take(&user, user_id).Error
	if err != nil {
		c.JSON(http.StatusOK, VideoList_response{
			Respond:   Respond{StatusCode: 1, StatusMsg: "账号不存在"},
			VideoList: []Video_respond{},
		})
	} else {
		l := len(user.Videos)
		Video_respondlist := make([]Video_respond, l)
		//倒序添加
		for i, x := range user.Videos {
			fmt.Println(int64(len(user.Videos)))
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
		c.JSON(http.StatusOK, VideoList_response{
			Respond:   Respond{StatusCode: 0},
			VideoList: Video_respondlist,
		})
	}
}
