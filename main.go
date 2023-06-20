package main

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

type MenuItem struct {
    Name      string `json:"name"`
    OrderCode string `json:"order_code"`
    Price     int    `json:"price"`
}

var foodMenu = []MenuItem{
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

func getFoodMenu(c echo.Context) error {
    return c.JSON(http.StatusOK, foodMenu)
}

func addFoodMenu(c echo.Context) error {
    item := new(MenuItem)
    if err := c.Bind(item); err != nil {
        return err
    }
    foodMenu = append(foodMenu, *item)
    return c.JSON(http.StatusOK, item)
}

func removeFoodMenu(c echo.Context) error {
    code := c.Param("orderCode")
    for i, item := range foodMenu {
        if item.OrderCode == code {
            foodMenu = append(foodMenu[:i], foodMenu[i+1:]...)
            return c.JSON(http.StatusOK, item)
        }
    }
    return c.NoContent(http.StatusNotFound)
}

func updateFoodMenu(c echo.Context) error {
    code := c.Param("orderCode")
    for i, item := range foodMenu {
        if item.OrderCode == code {
            newItem := new(MenuItem)
            if err := c.Bind(newItem); err != nil {
                return err
            }
            foodMenu[i] = *newItem
            return c.JSON(http.StatusOK, newItem)
        }
    }
    return c.NoContent(http.StatusNotFound)
}

func main() {
    e := echo.New()
    e.GET("/menu/food", getFoodMenu)
    e.POST("/menu/food", addFoodMenu)
    e.DELETE("/menu/food/:orderCode", removeFoodMenu)
    e.PUT("/menu/food/:orderCode", updateFoodMenu)
    e.Logger.Fatal(e.Start(":2000"))
}
