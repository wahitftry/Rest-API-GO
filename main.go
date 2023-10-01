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
    Nama        string `json:"nama"`
    KodePesanan string `json:"kode_pesanan"`
    Harga       int    `json:"harga"`
}

const (
    BatasMenuDefault     = 100
    UrutkanMenurutNama   = "nama"
    UrutkanMenurutHarga  = "harga"
    Menaik               = "naik"
    Menurun              = "turun"
    MaxHarga             = 1000000
    MaxPanjangNama       = 50
    MaxPanjangKodePesanan = 10
)

var menuMakanan = []MenuItem{
    {
        Nama:        "bakmie",
        KodePesanan: "bakmie",
        Harga:       12000,
    },
    {
        Nama:        "bakso",
        KodePesanan: "bakso",
        Harga:       8000,
    },
}

func validasiMenuItem(item *MenuItem) error {
    if item.Nama == "" || item.KodePesanan == "" {
        return errors.New("nama dan kode pesanan diperlukan")
    }
    if len(item.Nama) > MaxPanjangNama {
        return errors.New("nama terlalu panjang")
    }
    if len(item.KodePesanan) > MaxPanjangKodePesanan {
        return errors.New("kode pesanan terlalu panjang")
    }
    if item.Harga <= 0 || item.Harga > MaxHarga {
        return errors.New("harga tidak valid")
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

    e.PUT("/menu/:kode_pesanan", func(c echo.Context) error {
        kodePesanan := c.Param("kode_pesanan")

        var item *MenuItem
        for i := range menuMakanan {
            if menuMakanan[i].KodePesanan == kodePesanan {
                item = &menuMakanan[i]
                break
            }
        }

        if item == nil {
            return c.JSON(http.StatusNotFound, map[string]string{
                "pesan": "Menu tidak ditemukan",
            })
        }

        update := new(MenuItem)
        if err := c.Bind(update); err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{
                "pesan": "tidak dapat memproses permintaan",
            })
        }

        if update.Nama != "" {
            if len(update.Nama) > MaxPanjangNama {
                return c.JSON(http.StatusBadRequest, map[string]string{
                    "pesan": "nama terlalu panjang",
                })
            }
            item.Nama = update.Nama
        }

        if update.Harga != 0 {
            if update.Harga <= 0 || update.Harga > MaxHarga {
                return c.JSON(http.StatusBadRequest, map[string]string{
                    "pesan": "harga tidak valid",
                })
            }
            item.Harga = update.Harga
        }

        return c.JSON(http.StatusOK, map[string]interface{}{
            "pesan": "Menu berhasil diperbarui",
            "menu":  item,
        })
    })

    e.DELETE("/menu/:kode_pesanan", func(c echo.Context) error {
        kodePesanan := c.Param("kode_pesanan")

        for i := range menuMakanan {
            if menuMakanan[i].KodePesanan == kodePesanan {
                menuMakanan = append(menuMakanan[:i], menuMakanan[i+1:]...)
                return c.JSON(http.StatusOK, map[string]string{
                    "pesan": "Menu berhasil dihapus",
                })
            }
        }

        return c.JSON(http.StatusNotFound, map[string]string{
            "pesan": "Menu tidak ditemukan",
        })
    })

    e.Logger.Fatal(e.Start(":8080"))
}