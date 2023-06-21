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

const (
	DefaultMenuLimit = 100
	SortByNames      = "name"
	SortByPrice      = "price"
	Ascending        = "asc"
	Descending       = "desc"
)

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

func validateMenuItem(item *MenuItem) error {
	if item.Name == "" || item.OrderCode == "" {
		return errors.New("name and order code are required")
	}
	return nil
}

func validateLimit(limitStr string) (int, error) {
	if limitStr == "" {
		return DefaultMenuLimit, nil
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		return 0, errors.New("limit must be a positive integer")
	}
	return limit, nil
}

func filterMenu(menu []MenuItem, query string) []MenuItem {
	if query == "" {
		return menu
	}
	filteredMenu := make([]MenuItem, 0)
	for _, item := range menu {
		if strings.Contains(strings.ToLower(item.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(item.OrderCode), strings.ToLower(query)) {
			filteredMenu = append(filteredMenu, item)
		}
	}
	return filteredMenu
}

func (m MenuItem) lessThan(other MenuItem, sortBy string, order string) bool {
	switch sortBy {
	case SortByNames:
		if order == Descending {
			return m.Name > other.Name
		}
		return m.Name < other.Name
	case SortByPrice:
		if order == Descending {
			return m.Price > other.Price
		}
		return m.Price < other.Price
	default:
		return false
	}
}

type By func(m1, m2 *MenuItem) bool

func (by By) Sort(menu []MenuItem) {
	ms := &menuItemSorter{
		menu: menu,
		by:   by,
	}
	sort.Sort(ms)
}

type menuItemSorter struct {
	menu []MenuItem
	by   func(m1, m2 *MenuItem) bool
}

func (s *menuItemSorter) Len() int {
	return len(s.menu)
}

func (s *menuItemSorter) Swap(i, j int) {
	s.menu[i], s.menu[j] = s.menu[j], s.menu[i]
}

func (s *menuItemSorter) Less(i, j int) bool {
	return s.by(&s.menu[i], &s.menu[j])
}

func getFoodMenu(c echo.Context) error {
	q := c.QueryParam("q")
	sortBy := c.QueryParam("sort_by")
	order := c.QueryParam("order")
	limitStr := c.QueryParam("limit")
	limit, err := validateLimit(limitStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	filteredMenu := filterMenu(foodMenu, q)
	if sortBy != "" {
		By(func(m1, m2 *MenuItem) bool {
			return m1.lessThan(*m2, sortBy, order)
		}).Sort(filteredMenu)
	}
	if len(filteredMenu) > limit {
		filteredMenu = filteredMenu[:limit]
	}
	return c.JSON(http.StatusOK, filteredMenu)
}

func addFoodMenu(c echo.Context) error {
	item := new(MenuItem)
	if err := c.Bind(item); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := validateMenuItem(item); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
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
	for i := range foodMenu {
		if foodMenu[i].OrderCode == code {
			newItem := new(MenuItem)
			if err := c.Bind(newItem); err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
			if err := validateMenuItem(newItem); err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
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
