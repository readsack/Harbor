package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	_ "net/http/cookiejar"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client struct {
	http.Client
}

func main() {
	//fmt.Println("hello world")
	var c Client
	jar, _ := cookiejar.New(nil)
	c.Client.Jar = jar
	createAcc(&c, "hi", "hi", "hi", "http://localhost:8080/signup")
	logIntoAcc(&c, "hi", "hi", "http://localhost:8080/login")

	/*obj, _ := url.Parse("http://localhost:8080/")
	for _, cookie := range c.Client.Jar.Cookies(obj) {
		println(cookie.Value)
	}*/
	hey := make(chan string)
	defer close(hey)

	go func() {
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/msgs", nil)
		if err != nil {
			panic(err)
		}
		for {
			_, msg, _ := conn.ReadMessage()
			hey <- string(msg)
		}

	}()
	for {
		msgs := <-hey
		fmt.Println(msgs)
	}
	//createOrg(&c, "HelloWorldOrg", "http://localhost:8080/createorg")
	//createTeam(&c, "HelloTeam", "http://localhost:8080/createteam")
}

func createOrg(c *Client, name, url_str string) {
	resp, err := c.Client.PostForm(url_str, url.Values{
		"name": {name},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	println(string(body))
}

func inviteIntoOrg(c *Client, email, url_str string) {
	resp, err := c.Client.PostForm(url_str, url.Values{
		"email": {email},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	println(string(body))
}

func createTeam(c *Client, name, url_str string) {
	resp, err := c.Client.PostForm(url_str, url.Values{
		"name": {name},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	println(string(body))
}

func logIntoAcc(c *Client, email, pass, url_str string) {
	resp, err := c.Client.PostForm(url_str, url.Values{
		"email": {email},
		"pass":  {pass},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	println(string(body))
}

func createAcc(c *Client, name, email, pass, url_str string) (string, error) {
	resp, err := c.Client.PostForm(url_str, url.Values{
		"username": {name},
		"email":    {email},
		"pass":     {pass},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	println(string(body))
	return "", err
}
