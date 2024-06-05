package main

import (
	"html/template"
	"io"
	"log"
	"math"
	"net/mail"
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

type PageData struct {
	IsUserAuthenticated bool
	User                model.User
	Titles              []model.Title
	Title               model.Title
	SignupErrors        []string
	LoginErrors         []string
	RatingStars         [10]bool
}

func main() {
	if err := model.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	// FIXME: DO NOT STORE SESSION KEY IN SOURCE FILE!!!
	//        DO NOT USE THIS SESSION KEY IN PRODUCTION!!!
	model.InitUserSessions([]byte("6da35044863f15376abb9b27aa1c65dd01dfc31f98f1730ec4cfcb7f06ff10ba"))

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} error=${error}\n",
	}))
	e.Static("/static", "./static")

	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		user, isLogged, err := model.IsUserLoggedIn(c.Request(), c.Response().Writer)
		if err != nil {
			return c.String(500, err.Error())
		}

		titles, err := model.GetAllTitles()
		if err != nil {
			log.Fatal(err)
		}

		return c.Render(200, "index", PageData{Titles: titles, User: user, IsUserAuthenticated: isLogged})
	})

	e.GET("/title/:id", func(c echo.Context) error {
		user, isLogged, err := model.IsUserLoggedIn(c.Request(), c.Response().Writer)
		if err != nil {
			return c.String(500, err.Error())
		}

		titleIdStr := c.Param("id")
		titleId, err := strconv.Atoi(titleIdStr)
		if err != nil {
			return c.String(500, err.Error())
		}

		title, err := model.GetTitleById(uint64(titleId))
		if err != nil {
			return c.String(404, err.Error())
		}

		title.ConvertDuration()
		ratingStarsF, err := strconv.ParseFloat(title.RatingAvg, 32)
		if err != nil {
			return c.String(500, err.Error())
		}

		ratingStars := [10]bool{}
		ratingStarsNum := int(math.Round(ratingStarsF)) - 1
		if ratingStarsNum > 0 {
			ratingStars[ratingStarsNum] = true
		}

		return c.Render(200, "title", PageData{Title: title, User: user, IsUserAuthenticated: isLogged, RatingStars: ratingStars})
	})

	e.GET("/logout", func(c echo.Context) error {
		err := model.Logout(c.Request(), c.Response().Writer)
		if err != nil {
			c.String(502, err.Error())
		}

		return c.Redirect(302, "/")
	})

	e.GET("/signup", func(c echo.Context) error {
		_, isLogged, err := model.IsUserLoggedIn(c.Request(), c.Response().Writer)
		if err != nil {
			return c.String(500, err.Error())
		}

		if isLogged {
			return c.Redirect(302, "/")
		}

		return c.Render(200, "signup", nil)
	})

	e.POST("/signup", func(c echo.Context) error {
		_, isLogged, err := model.IsUserLoggedIn(c.Request(), c.Response().Writer)
		if err != nil {
			return c.String(500, err.Error())
		}

		if isLogged {
			return c.Redirect(302, "/")
		}

		action := c.FormValue("action")
		email := c.FormValue("email")
		username := c.FormValue("username")
		password := c.FormValue("password")

		if action == "login" {
			user, err := model.GetUserByEmail(email)
			if err != nil {
				return c.Render(200, "signup", PageData{LoginErrors: []string{"Неверный email или пароль"}})
			}

			if !user.CheckPassword(password) {
				return c.Render(200, "signup", PageData{LoginErrors: []string{"Неверный email или пароль"}})
			}

			if err := model.Login(c.Request(), c.Response().Writer, user.Id); err != nil {
				return c.String(500, err.Error())
			}

			return c.Redirect(302, "/")
		} else if action == "signup" {
			_, err := mail.ParseAddress(email)
			if err != nil {
				return c.Render(200, "signup", PageData{SignupErrors: []string{"Недействительный email"}})
			}

			emailExists, err := model.CheckUserEmailExists(email)
			if err != nil {
				c.Render(500, "signup", nil)
			}

			if emailExists {
				return c.Render(200, "signup", PageData{SignupErrors: []string{"Аккаунт с таким email уже существует"}})
			}

			userId, err := model.RegisterUser(email, username, password)
			if err != nil {
				if err.Error() == "Email already exists" {
					return c.Render(200, "signup", PageData{SignupErrors: []string{"Аккаунт с таким email уже существует"}})
				}

				return c.String(500, err.Error())
			}

			if err := model.Login(c.Request(), c.Response().Writer, userId); err != nil {
				return c.String(500, err.Error())
			}

			return c.Redirect(302, "/")
		}

		return c.String(500, "500")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
