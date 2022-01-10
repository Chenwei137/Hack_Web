package api

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type LoginJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Article struct {
	ArticleTitle    string `json:"title"`
	ArticleCategory string `json:"category"`
	ArticleText     string `json:"text"`
}

const (
	UserName string = ""
	Password string = ""
	Addr     string = "127.0.0.1"
	Port     int    = 3306
	Database string = "hack_account"
)

var (
	DB       *sql.DB
	username string
	password string
	status   int
)

func Login(logindata LoginJSON) int {
	var login_flag = 0
	var select_sql = "Select username,password,status From user_data Where username = '" + logindata.Username + "'"
	fmt.Print(select_sql)
	rows, err := DB.Query(select_sql)
	if err != nil {
		fmt.Println("1", err)
		//select data error
		login_flag = 1
	}
	total := 0
	for rows.Next() {
		err := rows.Scan(
			&username, &password, &status,
		)
		if err != nil {
			//search failed
			login_flag = 1
			fmt.Println("2", err)
			continue
		} else {
			//get search data
			total = 1
			if status == 1 {
				// status==1,account is active
				if logindata.Username == username && logindata.Password == password {
					login_flag = 200
				}
			} else {
				login_flag = 100
			}
		}
	}
	rows.Close()
	if total == 0 {
		//didn't have userdata
		login_flag = 10
	}
	return login_flag
}
func Create(account, password string) int {
	var create_flag = 0
	var select_sql = "select username from user_data Where username = '" + account + "'"
	rows, err := DB.Query(select_sql)
	if err != nil {
		//select fail
		fmt.Println(err)
		create_flag = 1
		return create_flag
	}
	total := 0
	for rows.Next() {
		total++
	}
	rows.Close()
	if total == 0 {
		var insert_sql = "Insert into user_data (username,password,status) Values (?,?,?)"
		_, err := DB.Exec(insert_sql, account, password, "1")
		//insert success
		if err != nil {
			create_flag = 10
			return create_flag
		}
		create_flag = 100
		return create_flag
	} else {
		//insert fail (account exist)
		create_flag = 1000
		return create_flag
	}
}

//func StoreArticle(user, title, category, text, date string)
func StoreArticle(title, category, text string) {
	articledata := []byte(text)
	t := time.Now()
	location, _ := time.LoadLocation("Asia/Taipei")
	local := strings.Split(t.In(location).String(), " ")
	localdate := local[0]
	localtime := strings.Split(local[1], ".")
	localdata := localdate + "-" + localtime[0]

	var articlepath = "/home/yuwei/go/src/GoWebServer/Article-storage/" + title + "-" + category + "-" + localdata + ".txt"
	f, err := os.Create(articlepath)
	if err != nil {
		fmt.Print(err)
	}
	_, err2 := f.Write(articledata)
	if err2 != nil {
		fmt.Print(err)
	}
	storedate := localdate + " " + localtime[0]
	var insert_sql = "Insert into file_data	 (username,title,category,filepath,storedate) Values (?,?,?,?,?)"
	_, err3 := DB.Exec(insert_sql, "None", title, category, articlepath, storedate)
	if err3 != nil {
		fmt.Print(err3)
	}
}
func ArticleTakeOut() {

}
func Con_database() *sql.DB {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, Password, Addr, Port, Database)
	var err error
	DB, err = sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("connection to mariadb failed:", err)
	} else {
		fmt.Println("conncetion to mariadb success", DB)
	}
	return DB
}
