package main

import (
	"html/template"
	"io"
	"log"
	"strconv"

	"github.com/Prost0i/kinoDB/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	if err := model.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "./static")

	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		type PageData struct {
			Titles []model.Title
		}

		titles, err := model.GetAllTitles()
		if err != nil {
			log.Fatal(err)
		}

		return c.Render(200, "index", PageData{Titles: titles})
	})

	e.GET("/title/:id", func(c echo.Context) error {
		titleIdStr := c.Param("id")
		titleId, err := strconv.Atoi(titleIdStr)
		if err != nil {
			log.Println(err)
			return c.Redirect(301, "/")
		}

		title, err := model.GetTitleById(uint64(titleId))
		if err != nil {
			log.Println(err)
			return c.Redirect(301, "/")
		}

		title.ConvertDuration()

		return c.Render(200, "title", title)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
