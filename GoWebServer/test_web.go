package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"test/api"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Claims struct {
	Account string `json:"account"`
	Role    string `json:"role"`
	jwt.StandardClaims
}

var jwtSecret = []byte("secret")

// web's html
func main_page(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", nil)
}
func login_page(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}
func create_page(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
}

func create_account(c *gin.Context) {
	var json_in api.LoginJSON
	err := c.ShouldBindJSON(&json_in)
	fmt.Println("check_get data")
	if err != nil {
		fmt.Print("check0")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		fmt.Println("check1")
		var create_result = api.Create(json_in.Username, json_in.Password)
		fmt.Println(create_result)
		fmt.Println("check2")
		if create_result == 1 {
			c.JSON(http.StatusOK, gin.H{
				"Welcome!": string(json_in.Username),
			})
			time.Sleep(1 * time.Second)
		} else if create_result == 0 {
			c.JSON(http.StatusOK, gin.H{
				"Welcome": "Create account fail! Please confirm whether the user exists!",
			})
			time.Sleep(1 * time.Second)
		}
	}
}
func article_accept(c *gin.Context) {
	var json_article api.Article
	err := c.ShouldBindJSON(&json_article)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		fmt.Print("標題: " + json_article.ArticleTitle + "\n")
		fmt.Print("類別: " + json_article.ArticleCategory + "\n")
		fmt.Print("內容: " + json_article.ArticleText + "\n")
		t := time.Now()
		location, _ := time.LoadLocation("Asia/Taipei")
		local := strings.Split(t.In(location).String(), " ")
		localdate := local[0]
		localtime := strings.Split(local[1], ".")
		localdata := localdate + "/" + localtime[0]
		fmt.Print(localdata)
		api.StoreArticle(json_article.ArticleTitle, json_article.ArticleCategory, json_article.ArticleText)
	}
}

/*func login(c *gin.Context){
	var(
		username string
		password string
	)
	if in, isExist := c.GetPostForm("username"); isExist && in != "" {
		username = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("必須輸入使用者名稱"),
		})
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		password = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("必須輸入密碼名稱"),
		})
		return
	}
	if err := Auth(username, password); err == nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"success": "登入成功",
		})
		return
	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": err,
		})
		return
	}
}*/
// mariadb's motion
func CreateTable(DB *sql.DB) {
	var sql = `CREATE TABLE IF NOT EXISTS user_data (
	username VARCHAR(64),
	password VARCHAR(64),
	status INT(4)
	); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create table failed:", err)
		return
	}
	fmt.Println("create table successd")
}

/*func InsertData(){
	result, err := db.Exec(
		"INSERT INTO user_info (name, age) VALUES (?, ?)",
		"syhlion",
		18,
	)
}*/

func main() {
	/*conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, Password, Addr, Port, Database)
	DB, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("connection to mariadb failed:", err)
		return
	}else {
		fmt.Println("conncetion to mariadb success", DB)
		CreateTable(DB)
	}*/
	api.Con_database()
	server := gin.Default()
	server.Static("/assets", "./template/assets")
	server.LoadHTMLGlob("template/html/*")
	server.GET("/login", login_page)
	server.GET("/create", create_page)
	server.GET("/main", main_page)
	server.POST("/login", func(c *gin.Context) {
		var json_in api.LoginJSON
		err := c.ShouldBindJSON(&json_in)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			var flag = api.Login(json_in)
			if flag == 1 {
				//sql error
				c.JSON(http.StatusUnauthorized, gin.H{"error": "SQL Failed"})
				fmt.Print("2")
			} else if flag == 10 {
				c.JSON(200, gin.H{"error": "User Empty"})
				return
			} else if flag == 100 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Password or Username was worng"})
			} else if flag == 200 {
				now := time.Now()
				jwtid := json_in.Username + strconv.FormatInt(now.Unix(), 10)
				role := "Member"

				claims := Claims{
					Account: json_in.Username,
					Role:    role,
					StandardClaims: jwt.StandardClaims{
						Audience:  json_in.Username,
						ExpiresAt: now.Add(5 * time.Minute).Unix(),
						Id:        jwtid,
						IssuedAt:  now.Unix(),
						Issuer:    "ginJWT",
						NotBefore: now.Add(10 * time.Second).Unix(),
						Subject:   json_in.Username,
					},
				}
				tokenClaims := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
				token, err := tokenClaims.SignedString(jwtSecret)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"token": token,
				})
				c.Request.URL.Path = "/main"
				server.HandleContext(c)
				return
			}
		}
		//c.JSON(200, gin.H{"status": json_in.Username})
	})
	server.POST("/acceptarticle", article_accept)
	server.POST("/create", create_account)
	server.Run(":8080")
}
