package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sinistra/go-gin-api/auth"
	"sinistra/go-gin-api/controllers"
	"strconv"
	"syscall"
	"time"

	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "sinistra/go-gin-api/docs" // docs is generated by Swag CLI, you have to import it.
)

var db *gorm.DB
var dbp *sqlx.DB
var mdb *mgo.Database
var port string
var host string

type (
	// todoModel describes a todoModel type
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}
	// transformedTodo represents a formatted todo
	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

func init() {
	godotenv.Load()
	var ok bool
	port, ok = os.LookupEnv("HOST_PORT")
	if !ok {
		port = "8080"
	}
	host, ok = os.LookupEnv("HOST_URL")
	if !ok {
		port = "http://localhost"
	}

	// open a db connection
	// var err error
	// db, err = gorm.Open("mysql", "root:root@/todos?charset=utf8&parseTime=True&loc=Local")
	// if err != nil {
	// 	panic("failed to connect database")
	// }
	// // Migrate the schema
	// db.AutoMigrate(&todoModel{})
}

// @title Swagger Todo API
// @version 1.0
// @description This is a sample Todo server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host todo.localhost
// @BasePath /v1
func main() {

	userController := controllers.UserController{}

	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// custom format
		return fmt.Sprintf("%s |%s|[%s]|%s|%s|%d|%s|%s|%s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.ClientIP,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "ResponseType"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./public", true)))

	urlString := fmt.Sprintf(host + ":" + port + "/swagger/doc.json")
	url := ginSwagger.URL(urlString) // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.POST("/add", addTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}

	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", showUser)

	router.POST("/login", userController.Login)
	router.GET("/checktoken", userController.CheckForToken)
	router.GET("/testauth", auth.JWTVerifyMiddleWare, userController.TestAuth)

	fmt.Println("Server is running at port " + port)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

// createTodo add a new todo
//
// @Summary Add a new todo
// @Description create todo based on JSON payload
// @Accept  json
// @Produce  json
// @Param   some_id     path    int     true        "Some ID"
// @Success 200 {object} string	"ok"
// @Failure 400 {string} string "400"
// @Failure 404 {string} string "404"
// @Router /api/v1/todo [post]
func createTodo(c *gin.Context) {
	// log.Println(c)
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{
		Title:     c.PostForm("title"),
		Completed: completed,
	}
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "Todo item created successfully!",
		"resourceId": todo.ID})
}

func addTodo(c *gin.Context) {
	// completed, _ := strconv.Atoi(c.PostForm("completed"))
	var todo todoModel
	if c.BindJSON(&todo) == nil {
		if todo.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Title cannot be empty", "data": todo})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "ok", "data": todo})
		}
	}
	// db.Save(&todo)
	// c.JSON(http.StatusCreated, gin.H{
	//     "status":     http.StatusCreated,
	//     "message":    "Todo item created successfully!",
	//     "resourceId": todo.ID})
}

//
// @Summary fetchAllTodo fetch all todos
// @Description retrieve all todos
// @Accept  json
// @Produce  json
// @Param   some_id     path    int     true        "Some ID"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "400"
// @Failure 404 {object} string "404"
// @Router /api/v1/todo [get]
func fetchAllTodo(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTodo
	db.Find(&todos)
	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No todo found!"})
		return
	}
	// transforms the todos for building a good response
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   _todos},
	)
}

// fetchSingleTodo fetch a single todo
func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}
	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

// updateTodo update a todo
func updateTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

// deleteTodo remove a todo
func deleteTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}

func showUser(c *gin.Context) {
	name := c.Param("name")
	page := c.DefaultQuery("page", "0")

	c.String(http.StatusOK, "Hello %s, page %s", name, page)
}
