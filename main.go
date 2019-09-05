package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var db *sql.DB
var world = []byte("world")

func init() {

	db, err := sql.Open("sqlite3", "lasers.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}

func dGet(key string) string {
	val := key
	return string(val)
}

func cartStore(val string) {

}

func clientStore(c echo.Context) error {
	name := c.FormValue("name")
	phone := c.FormValue("phone")
	comm := c.FormValue("comm")

	fmt.Println("Добавляем клиента:", name, phone, comm)
	return c.HTML(http.StatusOK, fmt.Sprintf("Клиент добавлен"))
}

func upload(c echo.Context) error {
	// Read form fields
	//name := c.FormValue("name")
	//znum := c.FormValue("znum")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create("storage/" + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<a href='/'>Следующий заказ</a>"))
}

func main() {

	e := echo.New()

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/upload", upload)
	e.POST("/clientAdd", clientStore)

	e.Logger.Fatal(e.Start(":1323"))
}
