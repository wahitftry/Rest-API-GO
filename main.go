package main

import (
	"github.com/labstack/echo/v4"
)

type MenuItem struct {
	Name      string
	OrderCode string
	price     int
}

func getFoodMenu(c echo.Context) error {
	foodMenu := []MenuItem{
		{
			Name:      "bakmie",
			OrderCode: "bakmie",
			price:     12000,
		},
		{
			Name:      "bakso",
			OrderCode: "bakso",
			price:     8000,
		},
	}
	return c.JSON(201, foodMenu)
}

func main() {
	e := echo.New()
	e.GET("menu/food", getFoodMenu)
	e.Logger.Fatal(e.Start(":2000"))
}
