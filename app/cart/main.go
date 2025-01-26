package main

import (
    "log"

    "TikTokMall/app/cart/handler"
    "TikTokMall/app/cart/kitex_gen/cart/cartservice"
)

func main() {
    svr := cartservice.NewServer(handler.NewCartServiceImpl())

    err := svr.Run()
    if err != nil {
        log.Println(err.Error())
    }
}
