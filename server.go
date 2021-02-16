package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)
var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
Addr:     "redis:6379",
Password: "", // no password set
DB:       0,  // use default DB
})

func TestGet(){
	
}

func InfoHandler(c echo.Context) error{
	link := c.Param("link")
	val, err := rdb.Get(ctx, link).Result()
	if err != nil {
		log.Println(err)
	}

	c.Redirect(http.StatusSeeOther,"https://" + val)
	return c.String(http.StatusOK, val)
}

func ShortifyHandler(c echo.Context) error {
	//fmt.Println("kuku")
	link := c.Param("link")
	fmt.Println(link)
	rand.Seed(time.Now().UnixNano())
	var b strings.Builder
	chars := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 8; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	err := rdb.Set(ctx,str, link, 0).Err()
	if err != nil {
		log.Println(err)
	}
	return c.String(http.StatusOK, "http://localhost:9000/"+str)
}

func main() {

	e := echo.New()
	
	e.POST("/shortify/:link",ShortifyHandler)
	e.GET("/:link",InfoHandler)

	e.Logger.Fatal(e.Start(":9000"))
}
