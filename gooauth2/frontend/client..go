package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
	router.GET("/auth/callback/google", loginCallback)
	router.GET("/auth/check/google", checkCodeState)
	router.Run(":8080")
}

/*func loginCallback(c *gin.Context) {
	platform := c.Param("platform")
	authInfo, err := getOAuth2Info(platform)
	if err != nil {
		log.Panic(err.Error())
		return
	}

	state := c.Request.FormValue("state")
	err = server.getState(state)
	if err != nil {
		log.Panic("invalid state.", err.Error())
		return
	}

	code := c.Request.FormValue("code")
	oauth2Token, err := authInfo.OAuth2Config.Exchange(c, code)
	if err != nil {
		log.Panic(authInfo.errs.exchangeFailed, err.Error())
		return
	}

	var reqURL string
	switch authInfo.platform {
	case "google":
		reqURL = authInfo.url + oauth2Token.AccessToken
	case "kakao":
		reqURL = authInfo.url
		authInfo.header.Set("Authorization", "Bearer "+oauth2Token.AccessToken)
	case "naver":
		reqURL = authInfo.url
		authInfo.header.Set("Authorization", "Bearer "+oauth2Token.AccessToken)
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Panic("oauth2 request failed. ", err.Error())
		return
	}

	req.Header = authInfo.header
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

	c.JSON(http.StatusOK, profile)
	//c.Data(http.StatusOK, "application/json", content)
	//fmt.Fprintf(ctx.Writer, "Response: %s", content)
}
*/

func loginCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func checkCodeState(c *gin.Context) {
	body := gin.H{
		"code": "",
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
}