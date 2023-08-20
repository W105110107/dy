package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Message_respond struct {
	MessageList []Mes `json:"message_list"` // 用户列表
	Respond     Respond
}

// Message
type Mes struct {
	Content    string `json:"content"`      // 消息内容
	CreateTime int64  `json:"create_time"`  // 消息发送时间 yyyy-MM-dd HH:MM:ss
	FromUserID int64  `json:"from_user_id"` // 消息发送者id
	ID         int64  `json:"id"`           // 消息id
	ToUserID   int64  `json:"to_user_id"`   // 消息接收者id
}

// 发送信息
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	content := c.Query("content")

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
		} else {
			err := DB.Transaction(func(tx *gorm.DB) error {
				if action_type == "1" {
					tuid, _ := strconv.ParseInt(to_user_id, 10, 64)
					newMessage := PeopleMessage{
						Con:        content,
						UserID:     tk.ID,
						TUID:       tuid,
						CreateTime: time.Now().Unix(),
					}
					if err := tx.Create(&newMessage).Error; err != nil {
						c.JSON(http.StatusOK, Respond{
							StatusCode: 1,
							StatusMsg:  "发送失败",
						})
						return err
					}
					c.JSON(http.StatusOK, Respond{StatusCode: 0})
					//这里似乎会自动调用一次  显示消息列表函数
					time.Sleep(time.Second)
				}
				return nil
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// 消息列表
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	pre_msg_time := c.Query("pre_msg_time")
	//不知道是不是客户端的bug，传回的时间戳位数过多，可能包含毫秒、
	if len(pre_msg_time) > 10 {
		pre_msg_time = pre_msg_time[:10]
	}
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
		} else {
			tuid, _ := strconv.ParseInt(to_user_id, 10, 64)
			p_m_t, _ := strconv.ParseInt(pre_msg_time, 10, 64)
			var user, touser User
			DB.Preload("PeopleMessages").Take(&user, tk.ID)
			DB.Preload("PeopleMessages").Take(&touser, tuid)
			List := append(user.PeopleMessages, touser.PeopleMessages...)
			sort.Sort(PM(List))
			var MessageList []Mes
			for _, x := range List {
				fmt.Println(x)
				if x.UserID == tk.ID && x.TUID == tuid || (x.UserID == tuid && x.TUID == tk.ID) {
					if x.CreateTime <= p_m_t {
						continue
					}
					MessageList = append(MessageList, Mes{
						Content:    x.Con,
						CreateTime: x.CreateTime,
						FromUserID: x.UserID,
						ID:         x.ID,
						ToUserID:   x.TUID,
					})
				}
			}
			c.JSON(http.StatusOK, Message_respond{
				Respond:     Respond{StatusCode: 0},
				MessageList: MessageList,
			})
		}
	}
}
