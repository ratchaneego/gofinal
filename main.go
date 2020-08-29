package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

var db *sql.DB

func init() {
	var err error
	//db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", "postgres://eihytncu:AJB4RN8B9nia4vAZwaqWAKJ3G36OyjB8@lallah.db.elephantsql.com:5432/eihytncu")
	if err != nil {
		log.Fatal(err)
	}
}

func createCustomersHandler(c *gin.Context) {
	t := Customer{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	row := db.QueryRow("INSERT INTO customers (id, name, email, status) values ($1, $2, $3) RETURNING id", t.Email, t.Email, t.Status)
	err := row.Scan(&t.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, t)
}

func getCustomersHandler(c *gin.Context) {
	status := c.Query("status")
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	customers := []Customer{}
	for rows.Next() {
		t := Customer{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		customers = append(customers, t)
	}
	tt := []Customer{}
	for _, item := range createCustomersHandler {
		if status != "" {
			if item.Status == status {
				tt = append(tt, item)
			}
		} else {
			tt = append(tt, item)
		}
	}
	c.JSON(http.StatusOK, tt)
}

func getCustomerByIdHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	row := stmt.QueryRow(id)
	t := &Customer{}
	err = row.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, t)
}

func authMiddleware(c *gin.Context) {
	fmt.Println("start #middleware")
	token := c.GetHeader("Authorization")
	if token != "Bearer token1234" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you don't have the permission!!"})
		c.Abort()
		return
	}
	c.Next()
	fmt.Println("end #middleware")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	apiV1.Use(authMiddleware)
	apiV1.GET("/customers", createCustomersHandler)
	apiV1.GET("/customers", getCustomersHandler)
	apiV1.GET("/customers/:id", getCustomerByIdHandler)
	return r
}

func main() {
	fmt.Println("customer service")
	r := setupRouter()
	r.Run(":2009")
}
