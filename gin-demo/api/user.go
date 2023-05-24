package api

import (
	"gin-demo/dao"
	"gin-demo/model"
	"gin-demo/utils"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

var rds redis.Conn

func register(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		utils.RespSuccess(c, "verification failed")
		return
	}
	Account := c.PostForm("account")
	Password := c.PostForm("password")
	Username := c.PostForm("username")

	flag := dao.SelectUser(Account)
	if flag {
		utils.RespFail(c, "user already exists")
		return
	}

	dao.AddUser(Account, Username, Password) //添加用户

	utils.RespSuccess(c, "add user successful")
}

func login(c *gin.Context) {
	if err := c.ShouldBind(&model.UserLogin{}); err != nil {
		utils.RespFail(c, "verification failed")
		return
	}
	account := c.PostForm("account")
	password := c.PostForm("password")

	flag := dao.SelectUser(account)
	if !flag {
		utils.RespFail(c, "user doesn't exists")
		return
	}

	selectPassword := dao.SelectPasswordFromAccount(account)

	if selectPassword != password {
		print(selectPassword)
		utils.RespFail(c, "wrong password")
		return
	}
	utils.RespSuccess(c, "login success")

}

func UserName(c *gin.Context) {
	rds = dao.RedisPOllInit().Get()
	username, _ := redis.String(rds.Do("lrange", "username", 0, -1))
	println(username)
	if len(username) <= 0 { //列表中没有
		println("缓存中查询不到")
		name := dao.Username()
		for _, p := range name {
			if p != "" {
				_, err := rds.Do("lpush", "username", p)
				if err != nil {
					println("错误")
				}
				rds.Do("expire", "username", 10)
			}
		}
	}
}
