package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type List_respond struct {
	Respond  Respond
	UserList []User_respond `json:"user_list"` // 用户信息列表
}

// 关注取关操作
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")

	if len(token) == 0 {
		c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "请先登录"})
	} else {
		var tk Token
		if err := DB.Where("Token = ?", token).Take(&tk).Error; err != nil {
			c.JSON(http.StatusOK, Respond{
				StatusCode: 1,
				StatusMsg:  "请先登录",
			})
		} else {
			err := DB.Transaction(func(tx *gorm.DB) error {
				touserid, _ := strconv.ParseInt(to_user_id, 10, 64)
				var touser, user User
				if err := tx.Take(&touser, touserid).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "关注失败",
					})
					return err
				}
				if err := tx.Take(&user, tk.ID).Error; err != nil {
					c.JSON(http.StatusOK, Respond{
						StatusCode: 1,
						StatusMsg:  "关注失败",
					})
					return err
				}
				if action_type == "1" {
					user.FollowList = append(user.FollowList, touserid)
					touser.FollowerList = append(touser.FollowerList, user.ID)
					if err := tx.Select("FollowList").Save(&user).Error; err != nil {
						c.JSON(http.StatusOK, Respond{
							StatusCode: 1,
							StatusMsg:  "关注失败",
						})
						return err
					}
					if err := tx.Select("FollowerList").Save(&touser).Error; err != nil {
						c.JSON(http.StatusOK, Respond{
							StatusCode: 1,
							StatusMsg:  "关注失败",
						})
						return err
					}
					c.JSON(http.StatusOK, Respond{
						StatusCode: 0,
						StatusMsg:  "关注成功",
					})
				} else {
					var index int
					index = -1
					//在用户的关注列表里删除作者的id
					for i, x := range user.FollowList {
						if x == touserid {
							index = i
							break
						}
					}
					if index == -1 {
						c.JSON(http.StatusOK, Respond{
							StatusCode: 1,
							StatusMsg:  "取消关注失败",
						})
						return errors.New("取消关注失败")
					}
					user.FollowList = append(user.FollowList[:index], user.FollowList[index+1:]...)
					//在作者的粉丝列表里删除用户的id
					index = -1
					for i, x := range touser.FollowerList {
						if x == user.ID {
							index = i
							break
						}
					}
					if index == -1 {
						c.JSON(http.StatusOK, Respond{
							StatusCode: 1,
							StatusMsg:  "取消关注失败",
						})
						return errors.New("取消关注失败")
					}
					touser.FollowerList = append(touser.FollowerList[:index], touser.FollowerList[index+1:]...)
					if err := tx.Select("FollowList").Save(&user).Error; err != nil {
						c.JSON(http.StatusOK, Respond{
							StatusCode: 1,
							StatusMsg:  "取消关注失败",
						})
						return err
					}
					if err := tx.Select("FollowerList").Save(&touser).Error; err != nil {
						c.JSON(http.StatusOK, Respond{
							StatusCode: 1,
							StatusMsg:  "取消关注失败",
						})
						return err
					}
					c.JSON(http.StatusOK, Respond{
						StatusCode: 0,
						StatusMsg:  "取关成功",
					})
				}
				return nil
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// 显示作者的关注列表
func FollowList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	if len(token) == 0 {
		c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "请先登录"})
	} else {
		var tk Token
		if err := DB.Where("Token = ?", token).Take(&tk).Error; err != nil {
			c.JSON(http.StatusOK, Respond{
				StatusCode: 1,
				StatusMsg:  "请先登录",
			})
		} else {
			var user User
			DB.Take(&user, user_id)
			l := len(user.FollowList)
			UserList := make([]User_respond, l)
			for i, x := range user.FollowList {
				//关注的用户的信息
				var tar User
				DB.Preload("Videos").Take(&tar, x)
				UserList[l-i-1] = User_respond{
					Avatar:          tar.Avatar,
					BackgroundImage: tar.BackgroundImage,
					FavoriteCount:   int64(len(tar.FavoriteList)),             // 喜欢数
					FollowCount:     int64(len(tar.FollowList)),               // 关注总数
					FollowerCount:   int64(len(tar.FollowerList)),             // 粉丝总数
					ID:              tar.ID,                                   // 用户id
					IsFollow:        getBool(&token, &(tar.FollowerList)),     // true-已关注，false-未关注
					Name:            tar.Name,                                 // 用户名称
					Signature:       tar.Signature,                            // 个人简介
					TotalFavorited:  strconv.FormatInt(tar.TotalFavorite, 10), // 获赞数量
					WorkCount:       int64(len(tar.Videos)),                   // 作品数
				}
			}
			c.JSON(http.StatusOK, List_respond{
				Respond:  Respond{StatusCode: 0},
				UserList: UserList,
			})
		}
	}
}

// 显示粉丝的关注列表
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	if len(token) == 0 {
		c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "请先登录"})
	} else {
		var tk Token
		if err := DB.Where("Token = ?", token).Take(&tk).Error; err != nil {
			c.JSON(http.StatusOK, Respond{
				StatusCode: 1,
				StatusMsg:  "请先登录",
			})
		} else {
			var user User
			DB.Take(&user, user_id)
			l := len(user.FollowerList)
			UserList := make([]User_respond, l)
			fmt.Println("测试：", l)
			for i, x := range user.FollowerList {
				//关注的用户的信息
				var tar User
				DB.Preload("Videos").Take(&tar, x)
				UserList[l-i-1] = User_respond{
					Avatar:          tar.Avatar,
					BackgroundImage: tar.BackgroundImage,
					FavoriteCount:   int64(len(tar.FavoriteList)),             // 喜欢数
					FollowCount:     int64(len(tar.FollowList)),               // 关注总数
					FollowerCount:   int64(len(tar.FollowerList)),             // 粉丝总数
					ID:              tar.ID,                                   // 用户id
					IsFollow:        getBool(&token, &(tar.FollowerList)),     // true-已关注，false-未关注
					Name:            tar.Name,                                 // 用户名称
					Signature:       tar.Signature,                            // 个人简介
					TotalFavorited:  strconv.FormatInt(tar.TotalFavorite, 10), // 获赞数量
					WorkCount:       int64(len(tar.Videos)),                   // 作品数
				}
			}
			c.JSON(http.StatusOK, List_respond{
				Respond:  Respond{StatusCode: 0},
				UserList: UserList,
			})
		}
	}
}

// 显示好友列表
func FriendList(c *gin.Context) {
	token := c.Query("token")
	if len(token) == 0 {
		c.JSON(http.StatusOK, Respond{StatusCode: 1, StatusMsg: "请先登录"})
	} else {
		var tk Token
		if err := DB.Where("Token = ?", token).Take(&tk).Error; err != nil {
			c.JSON(http.StatusOK, Respond{
				StatusCode: 1,
				StatusMsg:  "请先登录",
			})
		} else {
			var user User
			DB.Take(&user, tk.ID)
			l := len(user.FollowList)
			UserList := make([]User_respond, l)
			for i, x := range user.FollowList {
				//关注的用户的信息
				var tar User
				DB.Preload("Videos").Take(&tar, x)
				UserList[l-i-1] = User_respond{
					Avatar:          tar.Avatar,
					BackgroundImage: tar.BackgroundImage,
					FavoriteCount:   int64(len(tar.FavoriteList)),             // 喜欢数
					FollowCount:     int64(len(tar.FollowList)),               // 关注总数
					FollowerCount:   int64(len(tar.FollowerList)),             // 粉丝总数
					ID:              tar.ID,                                   // 用户id
					IsFollow:        getBool(&token, &(tar.FollowerList)),     // true-已关注，false-未关注
					Name:            tar.Name,                                 // 用户名称
					Signature:       tar.Signature,                            // 个人简介
					TotalFavorited:  strconv.FormatInt(tar.TotalFavorite, 10), // 获赞数量
					WorkCount:       int64(len(tar.Videos)),                   // 作品数
				}
			}
			c.JSON(http.StatusOK, List_respond{
				Respond:  Respond{StatusCode: 0},
				UserList: UserList,
			})
		}
	}
}
