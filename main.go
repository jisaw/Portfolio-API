package main

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
	"strconv"
)

type Article struct {
	Id int64 `db:"article_id"`
	Created int64
	Title string
	Content string
}

type Contact struct {
	Id int64 `db:"contact_id"`
	Created int64
	Name string
	Title string
	Company string
	Email string
	Message string
	Phone int64
}

type Login struct {
	Id int64 `db:"login_id"`
	Created int64
	Username string
	Password string
	Create bool
}

var dbmap = initDb()

func initDb() gorp.DbMap {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	checkErr(err, "sql.Open faild")
	dbmap := gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Article{}, "articles").SetKeys(true, "Id")
	dbmap.AddTableWithName(Contact{}, "contacts").SetKeys(true, "Id")
	dbmap.AddTableWithName(Login{}, "logins").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Printf(msg, err)//.(*errors.Error).ErrorStack())
	}
}

func index (c *gin.Context) {
	content := gin.H{"Hello": "World"}
	c.JSON(200, content)
}

func ArticlesList(c *gin.Context) {
	var articles []Article
	_, err := dbmap.Select(&articles, "select * from articles order by article_id")
	checkErr(err, "Select failed")
	content := gin.H{"records": articles,}
	c.JSON(200, content)
}

func ArticlesDetail(c *gin.Context) {
	article_id := c.Params.ByName("id")
	a_id, _ := strconv.Atoi(article_id)
	article := getArticle(a_id)
	content := gin.H{"title": article.Title, "content": article.Content}
	c.JSON(200, content)
}

func ArticlePost(c *gin.Context) {
	var json Article

	c.Bind(&json)
	article := createArticle(json.Title, json.Content)
	if article.Title == json.Title {
		content := gin.H{
			"result": "Success",
			"title": article.Title,
			"content": article.Content,
		}
		c.JSON(201, content)
	} else {
		c.JSON(500, gin.H{"result": "An error occured"})
	}
}

func createArticle(title, body string) Article {
	article := Article{
		Created: time.Now().UnixNano(),
		Title: title,
		Content: body,
	}

	err := dbmap.Insert(&article)
	checkErr(err, "Insert failed")
	return article
}

func getArticle(article_id int) Article {
	article := Article{}
	err := dbmap.SelectOne(&article, "select * from articles where article_id=?", article_id)
	checkErr(err, "selectOne failed")
	return article
}

func ContactsList(c *gin.Context) {
	var contacts []Contact
	_, err := dbmap.Select(&contacts, "select * from contacts order by contact_id")
	checkErr(err, "Select Failed")
	content := gin.H{"records": contacts,}
	c.JSON(200, content)
}

func ContactPost(c *gin.Context) {
	var json Contact

	c.Bind(&json)
	contact := createContact(json.Name, json.Title, json.Company, json.Email, json.Message, json.Phone)
	if contact.Title == json.Title {
		content := gin.H{
			"result": "Success",
			"name": contact.Name,
			"title": contact.Title,
			"comapany": contact.Company,
			"email": contact.Email,
			"message": contact.Message,
			"phone": contact.Phone,
		}
		c.JSON(201, content)
	} else {
		c.JSON(500, gin.H{"result": "An error occured"})
	}
}

func ContactsDetail(c *gin.Context) {
	contact_id := c.Params.ByName("id")
	c_id, _ := strconv.Atoi(contact_id)
	contact := getContact(c_id)
	content := gin.H{
		"name": contact.Name,
		"title": contact.Title,
		"comapany": contact.Company,
		"email": contact.Email,
		"message": contact.Message,
		"phone": contact.Phone,
	}
	c.JSON(200, content)
}

func createContact(name, title, company, email, message string, phone int64) Contact {
	contact := Contact{
		Created: time.Now().UnixNano(),
		Name: name,
		Title: title,
		Company: company,
		Email: email,
		Message: message,
		Phone: phone,
	}

	err := dbmap.Insert(&contact)
	checkErr(err, "Insert Failed")
	return contact
}

func getContact(contact_id int) Contact {
	contact := Contact{}
	err := dbmap.SelectOne(&contact, "select * from contacts where contact_id=?", contact_id)
	checkErr(err, "SelectOne Failed")
	return contact
}

func LoginPost(c *gin.Context) {
	var json Login

	c.Bind(&json)
	if json.Create == true{
		createLogin(json.Username, json.Password)
	} else {
		login := checkLogin(json.Username, json.Password)
		if login == true {
			c.JSON(200, gin.H{"result": "success"})
			} else {
				c.JSON(500, gin.H{"result": "error"})
			}
		}
}

func checkLogin(username, password string) bool {
	auth := Login{}
	err := dbmap.SelectOne(&auth, "select * from logins where username=?", username)
	checkErr(err, "selectOne Failed")
	if password == auth.Password {
		return true
	} else {
		return false
	}
}

func createLogin(username, password string) Login {
	login := Login{
		Created: time.Now().UnixNano(),
		Username: username,
		Password: password,
	}

	err := dbmap.Insert(&login)
	checkErr(err, "Insert Failed")
	return login
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept=Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-with")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	}
}

func main() {
	defer dbmap.Db.Close()

	app := gin.Default()
	app.Use(CORSMiddleware())
	app.GET("/", index)
	app.GET("/articles", ArticlesList)
	app.POST("/articles", ArticlePost)
	app.GET("/articles/:id", ArticlesDetail)

	app.GET("/contacts", ContactsList)
	app.POST("/contacts", ContactPost)
	app.GET("/contacts/:id", ContactsDetail)

	app.POST("/login", LoginPost)

	app.Run(":8000")
}
