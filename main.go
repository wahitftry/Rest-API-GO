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
    Nama      string `json:"nama"`
    KodePesanan string `json:"kode_pesanan"`
    Harga     int    `json:"harga"`
}

const (
    BatasMenuDefault = 100
    UrutkanMenurutNama = "nama"
    UrutkanMenurutHarga = "harga"
    Menaik = "naik"
    Menurun = "turun"
)

var menuMakanan = []MenuItem{
    {
        Nama:      "bakmie",
        KodePesanan: "bakmie",
        Harga:     12000,
    },
    {
        Nama:      "bakso",
        KodePesanan: "bakso",
        Harga:     8000,
    },
}

func validasiMenuItem(item *MenuItem) error {
    if item.Nama == "" || item.KodePesanan == "" {
        return errors.New("nama dan kode pesanan diperlukan")
    }
    return nil
}

func validasiBatas(batasStr string) (int, error) {
    if batasStr == "" {
        return BatasMenuDefault, nil
    }
    batas, err := strconv.Atoi(batasStr)
    if err != nil {
        return 0, errors.New("batas harus berupa angka")
    }
    if batas <= 0 {
        return 0, errors.New("batas harus lebih besar dari 0")
    }
    return batas, nil
}

func urutkanMenu(menu []MenuItem, urutan string, arah string) []MenuItem {
    switch urutan {
    case UrutkanMenurutNama:
        if arah == Menaik {
            sort.Slice(menu, func(i, j int) bool {
                return strings.ToLower(menu[i].Nama) < strings.ToLower(menu[j].Nama)
            })
        } else {
            sort.Slice(menu, func(i, j int) bool {
                return strings.ToLower(menu[i].Nama) > strings.ToLower(menu[j].Nama)
            })
        }
    case UrutkanMenurutHarga:
        if arah == Menaik {
            sort.Slice(menu, func(i, j int) bool {
                return menu[i].Harga < menu[j].Harga
            })
        } else {
            sort.Slice(menu, func(i, j int) bool {
                return menu[i].Harga > menu[j].Harga
            })
        }
    }
    return menu
}

func main() {
    e := echo.New()

    e.GET("/menu", func(c echo.Context) error {
        batasStr := c.QueryParam("batas")
        urutan := c.QueryParam("urutan")
        arah := c.QueryParam("arah")

        batas, err := validasiBatas(batasStr)
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{
                "pesan": err.Error(),
            })
        }

        menu := menuMakanan[:batas]
        menu = urutkanMenu(menu, urutan, arah)

        return c.JSON(http.StatusOK, map[string]interface{}{
            "pesan": "Berikut adalah menu makanan",
            "menu":  menu,
        })
    })

    e.POST("/menu", func(c echo.Context) error {
        item := new(MenuItem)
        if err := c.Bind(item); err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{
                "pesan": "tidak dapat memproses permintaan",
            })
        }

        if err := validasiMenuItem(item); err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{
                "pesan": err.Error(),
            })
        }

        menuMakanan = append(menuMakanan, *item)

        return c.JSON(http.StatusCreated, map[string]interface{}{
            "pesan": "Menu berhasil ditambahkan",
            "menu":  item,
        })
    })

    e.Logger.Fatal(e.Start(":8080"))
}