package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("frontend/templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
	router.GET("/result", func(c *gin.Context) {
		email, err := c.Cookie("email")
		if err != nil {
			log.Println("get email cookie failed. ", err)
		}
		id, err := c.Cookie("id")
		if err != nil {
			log.Println("get id cookie failed. ", err)
		}
		picture, err := c.Cookie("picture")
		if err != nil {
			log.Println("get picture cookie failed. ", err)
		}
		access_token, err := c.Cookie("access_token")
		if err != nil {
			log.Println("get access token cookie failed. ", err)
		}
		refresh_token, err := c.Cookie("refresh_token")
		if err != nil {
			log.Println("get refresh token cookie failed. ", err)
		}
		c.HTML(http.StatusOK, "result.tmpl", gin.H{
			"email":         email,
			"id":            id,
			"picture":       picture,
			"access_token":  access_token,
			"refresh_token": refresh_token,
		})
	})
	router.GET("/auth/callback/google", loginCallback)
	router.GET("/auth/check/google", checkCodeState)
	router.Run(":8080")
}

func loginCallback(c *gin.Context) {
	state := c.Request.FormValue("state")
	code := c.Request.FormValue("code")
	body := gin.H{
		"state": state,
		"code":  code,
	}
	data, err := json.Marshal(body)
	if err != nil {
		log.Fatal("marshal failed. ", err)
	}
	reqURL := fmt.Sprintf("http://localhost:8090/auth/callback/google?state=%s&code=%s", state, code)
	req, err := http.NewRequest(http.MethodGet, reqURL, bytes.NewReader(data))
	if err != nil {
		log.Fatal("request failed. ", err)
	}
	client := http.Client{}
	resp, err := client.Do(req)

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic("read response body failed. ", err.Error())
		return
	}
	defer func(closer io.ReadCloser) {
		err := closer.Close()
		if err != nil {
			log.Panic("oauth2 response.body close failed. ", err.Error())
			return
		}
	}(resp.Body)

	var profile map[string]interface{}
	err = json.Unmarshal(content, &profile)
	if err != nil {
		log.Panic("oauth2 profile unmarshal failed. ", err.Error())
		return
	}

	c.SetCookie("access_token", profile["access_token"].(string), 10, "/", c.Request.URL.Hostname(), false, true)
	c.SetCookie("refresh_token", profile["refresh_token"].(string), 10, "/", c.Request.URL.Hostname(), false, true)
	c.SetCookie("email", profile["email"].(string), 10, "/", c.Request.URL.Hostname(), false, true)
	c.SetCookie("id", profile["id"].(string), 10, "/", c.Request.URL.Hostname(), false, true)
	c.SetCookie("picture", profile["picture"].(string), 10, "/", c.Request.URL.Hostname(), false, true)
	location := url.URL{
		Path: "/result",
	}
	c.Redirect(http.StatusTemporaryRedirect, location.RequestURI())
	//c.JSON(http.StatusOK, profile)
}

func checkCodeState(c *gin.Context) {
	body := gin.H{
		"code":  "",
		"state": "",
	}
	data, err := json.Marshal(body)
	if err != nil {
		log.Fatal("body marshal failed. ", err)
	}

	client := http.Client{}
	reqURL := "http://localhost:8090/auth/check/google"
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(data))
	if err != nil {
		return
	}
	log.Println(req)

	//req.Header = authInfo.header
	resp, err := client.Do(req)
	if err != nil {
		log.Panic("oauth2 response failed. ", err.Error())
		return
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic("read response body failed. ", err.Error())
		return
	}
	defer func(closer io.ReadCloser) {
		err := closer.Close()
		if err != nil {
			log.Panic("oauth2 response.body close failed. ", err.Error())
			return
		}
	}(resp.Body)

	var profile map[string]interface{}
	err = json.Unmarshal(content, &profile)
	if err != nil {
		log.Panic("oauth2 profile unmarshal failed. ", err.Error())
		return
	}
	log.Println(profile)

	c.JSON(http.StatusOK, profile)
}
