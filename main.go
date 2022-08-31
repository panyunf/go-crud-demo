package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

type User struct {
	gorm.Model        //给结构体添加主键
	Name       string `gorm:"type:varchar(20);not null"`
	Password   string `gorm:"type:varchar(20);not null"`
}

func main() {
	db := InitDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	r := gin.Default()

	//注册
	r.POST("/register", func(c *gin.Context) {

		//获取参数
		name := c.PostForm("name")
		password := c.PostForm("password")

		//数据验证
		if len(name) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名不能为空",
			})
			return
		}
		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码不能少于六位",
			})
		}
		var user User
		db.Where("name = ?", "name").First(&user)
		if user.ID != 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户已存在",
			})
			return

			//创建用户
			hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"code":    500,
					"message": "密码加密错误",
				})
				return
			}
			newUser := User{
				Name:     name,
				Password: string(hasedPassword),
			}
			db.Create(&newUser)
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "注册成功",
			})

		}

	})
	//登录
	r.POST("/login", func(c *gin.Context) {
		Password := c.PostForm("Password")
		if len(Password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码不能少于六位",
			})
		}
		var user User
		db.Where("name = ?").First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户不存在",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password)); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码错误",
			})

		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登录成功",
		})
		//测试
		//r.GET("/", func(c *gin.Context) {
		//	c.JSON(200, gin.H{
		//		"message": "请求成功",
		//	})
		//})

		//增加
		//r.POST("", func(c *gin.Context) {
		//	var data list
		//	err := c.ShouldBindJSON(&data)
		//	//判断绑定是否成功
		//	if err != nil {
		//		c.JSON(200, gin.H{
		//			"msg":  "添加失败",
		//			"data": gin.H{},
		//			"code": 400,
		//		})
		//
		//	} else {
		//
		//		//数据库的操作
		//		db.Create(&data) //创建一条数据
		//
		//		c.JSON(200, gin.H{
		//			"msg":  "添加成功",
		//			"data": data,
		//			"code": 200,
		//		})
		//	}
		//})

	})
	r.Run(":3040")

}
func InitDB() *gorm.DB {
	dsn := "root:117351@tcp(127.0.0.1:3306)/crud-list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	fmt.Println(db)
	fmt.Println(err)

	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}
