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

func dGet(key string) string {
	val := key
	return string(val)
}

func cartStore(val string) {

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

	db, err := sql.Open("sqlite3", "lasers.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e := echo.New()

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/", func(e echo.Context) error {
		qtype := e.FormValue("qtype")

		if qtype == "addclient" {
			fmt.Println("ADD CLIENT MODE")
		}

		return e.String("301", "/")
	})
	e.POST("/upload", upload)
	e.POST("/clientAdd", func(e echo.Context) error {
		name := e.FormValue("name")
		phone := e.FormValue("phone")
		comm := e.FormValue("comm")

		fmt.Println("Добавляем клиента:", name, phone, comm)
		return e.String(200, "OK")
	})

	e.POST("/getClient", func(e echo.Context) error {
		var name string
		var phone string
		clcl := e.FormValue("clid")
		err := db.QueryRow("SELECT name, phone FROM clients WHERE id=?", clcl).Scan(&name, &phone)
		if err != nil {
			log.Fatal(err)
		}
		ret := name + ", tel: " + phone
		return e.String(200, ret)
	})

	e.POST("/clist", func(e echo.Context) error {
		type Client struct {
			ID      string
			name    string
			phone   string
			comment string
		}

		ret := "<option value='0'>Дед Нихто</option>"

		rows, err := db.Query("SELECT * FROM clients")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		cls := make([]*Client, 0)
		for rows.Next() {
			cl := new(Client)
			err := rows.Scan(&cl.ID, &cl.name, &cl.phone, &cl.comment)
			if err != nil {
				log.Fatal(err)
			}
			cls = append(cls, cl)
		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		for _, cl := range cls {
			fmt.Printf("%s, %s, %s, %s\n", cl.ID, cl.name, cl.phone, cl.comment)
			ret = ret + "<option value='" + cl.ID + "'>" + cl.name + ", tel: " + cl.phone + "</option>"
		}
		return e.String(200, ret)
	})

	e.Logger.Fatal(e.Start(":80"))
}
