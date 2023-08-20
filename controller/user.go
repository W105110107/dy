package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserRegister_respond struct {
	Respond Respond
	Token   string `json:"token"`   // 用户鉴权token
	UserID  int64  `json:"user_id"` // 用户id
}

type UserInfo_respond struct {
	Respond  Respond
	UserInfo User_respond `json:"user"` // 用户信息
}

// 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	// 先判空
	if len(username) == 0 || len(password) == 0 {
		if len(username) == 0 {
			c.JSON(http.StatusOK, UserRegister_respond{
				Respond: Respond{StatusCode: 1, StatusMsg: "账号不可为空"},
				Token:   "",
				UserID:  -1,
			})
		} else {
			c.JSON(http.StatusOK, UserRegister_respond{
				Respond: Respond{StatusCode: 1, StatusMsg: "密码不可为空"},
				Token:   "",
				UserID:  -1,
			})
		}
	}
	// 校验参数长度
	if len(password) > 32 || len(password) <= 5 || len(username) > 32 {
		if len(username) > 32 {
			c.JSON(http.StatusOK, UserRegister_respond{
				Respond: Respond{StatusCode: 1, StatusMsg: "账号长度应不超过32位"},
				Token:   "",
				UserID:  -1,
			})
		} else {
			c.JSON(http.StatusOK, UserRegister_respond{
				Respond: Respond{StatusCode: 1, StatusMsg: "密码长度应不超过32位且不小于6位"},
				Token:   "",
				UserID:  -1,
			})
		}
	}
	token := username + "-" + password
	var tk Token
	err := DB.Where("Token = ?", token).Take(&tk).Error
	if err == nil {
		//用户已存在
		c.JSON(http.StatusOK, UserRegister_respond{
			Respond: Respond{StatusCode: 1, StatusMsg: "用户已存在"},
			Token:   token,
			UserID:  tk.ID,
		})
	} else {
		//新建用户
		user := User{
			Avatar:          DUP,
			BackgroundImage: DUHI,
			Signature:       username + "的个人简介",
			TotalFavorite:   0,
			Name:            username,
			FollowList:      []int64{},
			FollowerList:    []int64{},
			FavoriteList:    []int64{},
			Videos:          []Video{},
		}
		DB.Create(&user)
		newToken := Token{
			Token: token,
			ID:    user.ID,
		}
		DB.Create(newToken)
		c.JSON(http.StatusOK, UserRegister_respond{
			Respond: Respond{StatusCode: 0, StatusMsg: "创建成功"},
			Token:   token,
			UserID:  user.ID,
		})
	}
}

// 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + "-" + password
	var tk Token
	err := DB.Where("Token = ?", token).Take(&tk).Error
	if err == nil {
		c.JSON(http.StatusOK, UserRegister_respond{
			Respond: Respond{StatusCode: 0, StatusMsg: "登录成功"},
			Token:   token,
			UserID:  tk.ID,
		})
	} else {
		c.JSON(http.StatusOK, UserRegister_respond{
			Respond: Respond{StatusCode: 1, StatusMsg: "账号不存在"},
			Token:   "",
			UserID:  -1,
		})
	}
}

// 用户信息
func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	var user User
	err := DB.Preload("Videos").Take(&user, user_id).Error
	if err != nil {
		c.JSON(http.StatusOK, UserInfo_respond{
			Respond: Respond{StatusCode: 1, StatusMsg: "账号不存在"},
		})
	} else {
		c.JSON(http.StatusOK, UserInfo_respond{
			Respond: Respond{StatusCode: 0},
			UserInfo: User_respond{
				Avatar:          user.Avatar,
				BackgroundImage: user.BackgroundImage,
				FavoriteCount:   int64(len(user.FavoriteList)),
				FollowCount:     int64(len(user.FollowList)),
				FollowerCount:   int64(len(user.FollowerList)),
				ID:              user.ID,
				IsFollow:        getBool(&token, &user.FollowerList),
				Name:            user.Name,
				Signature:       user.Signature,
				TotalFavorited:  strconv.FormatInt(user.TotalFavorite, 10),
				WorkCount:       int64(len(user.Videos)),
			},
		})
	}
}
