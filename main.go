package main

import (
    "errors"
    "net/http"
    "sort"
    "strconv"
    "strings"

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
    q := c.QueryParam("q")
    sortBy := c.QueryParam("sort_by")
    order := c.QueryParam("order")
    limit, err := strconv.Atoi(c.QueryParam("limit"))
    if err != nil {
        limit = len(foodMenu)
    }
    filteredMenu := foodMenu
    if q != "" {
        filteredMenu = make([]MenuItem, 0)
        for _, item := range foodMenu {
            if strings.Contains(strings.ToLower(item.Name), strings.ToLower(q)) ||
                strings.Contains(strings.ToLower(item.OrderCode), strings.ToLower(q)) {
                filteredMenu = append(filteredMenu, item)
            }
        }
    }
    if sortBy != "" {
        sort.Slice(filteredMenu, func(i, j int) bool {
            switch sortBy {
            case "name":
                if order == "desc" {
                    return filteredMenu[i].Name > filteredMenu[j].Name
                }
                return filteredMenu[i].Name < filteredMenu[j].Name
            case "price":
                if order == "desc" {
                    return filteredMenu[i].Price > filteredMenu[j].Price
                }
                return filteredMenu[i].Price < filteredMenu[j].Price
            default:
                return i < j
            }
        })
    }
    if len(filteredMenu) > limit {
        filteredMenu = filteredMenu[:limit]
    }

    return c.JSON(http.StatusOK, filteredMenu)
}

func addFoodMenu(c echo.Context) error {
    item := new(MenuItem)
    if err := c.Bind(item); err != nil {
        return err
    }
    if item.Name == "" || item.OrderCode == "" {
        return errors.New("Invalid input")
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
            if newItem.Name == "" || newItem.OrderCode == "" {
                return errors.New("Invalid input")
            }
            foodMenu[i] = *newItem
            return c.JSON(http.StatusOK, newItem)
        }
    }
    return c.NoContent(http.StatusNotFound)
}

func main() {
    e := echo.New()

    e.GET("/menu", getFoodMenu)
    e.POST("/menu", addFoodMenu)
    e.DELETE("/menu/:orderCode", removeFoodMenu)
    e.PUT("/menu/:orderCode", updateFoodMenu)

    e.Logger.Fatal(e.Start(":8080"))
}
