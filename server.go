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

	"io/ioutil"
)

type site struct {
    name string
    code int
}

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
Addr:     "redis:6379",
Password: "", 
DB:       0,  
})

func testsite(site1 string, ch chan *site){
	buf := site{name: site1, code:404}
	resp, err := http.Get("https://" + site1)
	if err != nil {
		buf = site{name: site1, code:404}
	}else
	{ 	_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			buf = site{name: site1, code:404}
		}else{
			buf = site{name: site1, code:200}
		}
	}	
	ch <- &buf
}

func TestGet(){

	for {
		val := rdb.Do(ctx, "KEYS", "*").String()
		words := strings.Fields(val)
		words[0] = ""
		words[1] = ""
		words[2] = words[2][1:]
		words[len(words) -1] = words[len(words) -1 ][:len(words[len(words) -1 ]) -1]
		ch := make(chan *site,len(words))
		my := make([]site,len(words))
		for _, res := range words{
			if res != ""{
				val, err := rdb.Get(ctx, res).Result()
				if err != nil {
					log.Println(err)
				}
				go testsite(val,ch)
			}
		}
		for j:= 0; j < len(words) - 2;j++{
			my[j]= *(<-ch)
			fmt.Println(my[j])
		}
		time.Sleep(600 * time.Second)
	}
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
	go TestGet()
	e.Logger.Fatal(e.Start(":9000"))
}
