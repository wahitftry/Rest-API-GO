package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// mendefinisikan struk ini tipe data
type MenuItem struct {
	Name      string
	OrderCode string
	Price     int
}

func getFoodMenu(c echo.Context) error {
	foodMenu := []MenuItem{
		{
			Name:      "bakmie",
			OrderCode: "bakmie",
			Price:     12000,
		},
		{
			Name:      "bakso",
			OrderCode: "bakso",
			Price:     8000,
		},
	}
	return c.JSON(http.StatusOK, foodMenu) // menjalankan status kode JSON
}

func main() {
	e := echo.New()
	e.GET("menu/food", getFoodMenu)
	e.Logger.Fatal(e.Start(":2000"))
}
