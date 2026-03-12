package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func main() {
	//fmt.Println("hello world")
	c := new(http.Client)
	jar, _ := cookiejar.New(nil)
	c.Jar = jar
	resp, err := c.PostForm("http://localhost:8080/login", url.Values{
		"email": {"hello"},
		"pass":  {"passwd"},
	})
	if err != nil {
		log.Fatal(err)
	}
	ck, _ := http.ParseSetCookie(resp.Header["Set-Cookie"][0])
	fmt.Println(ck.Value)
	urlObj, err := url.Parse("http://localhost:8080/")
	if err != nil {
		log.Println(err)
	}
	c.Jar.SetCookies(urlObj, []*http.Cookie{ck})

}
