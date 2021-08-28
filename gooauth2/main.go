package main

import (
	"context"
	"fmt"
	"gooauth2/api"
	"gooauth2/config"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

var ctx = context.Background()

// var (
// 	router = gin.Default()
// )

func main() {
	//
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config.")
	}

	server, err := api.NewServer(cfg)
	if err != nil {
		log.Fatal("cannot create new server.")
	}

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server.")
	}

	// router.POST("/login", server.Login)
	// router.POST("/todo", TokenAuthMiddleware(), CreateTodo)
	// router.POST("/logout", TokenAuthMiddleware(), Logout)
	// router.POST("/token/refresh", Refresh)
	// router.GET("/all", TokenAuthMiddleware(), func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, "OK")
	// })

	// //clientID := os.Getenv("GOOGLE_CLIENT_ID")
	// //clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	// //redirectURL := os.Getenv("GOOGLE_REDIRECT_URI")

	// oauth2Conf := &oauth2.Config{
	// 	ClientID:     "<CLIENT ID>",
	// 	ClientSecret: "<CLIENT SECRET>",
	// 	Endpoint:     google.Endpoint,
	// 	RedirectURL:  "<CALLBACK URL>",
	// 	Scopes: []string{
	// 		"https://www.googleapis.com/auth/userinfo.email",
	// 		//"https://www.googleapis.com/auth/userinfo.profile",
	// 	},
	// }
	// router.GET("/auth/google", func(ctx *gin.Context) {
	// 	url := oauth2Conf.AuthCodeURL("random string")
	// 	ctx.Redirect(http.StatusTemporaryRedirect, url)
	// })
	// router.GET("/auth/callback", func(ctx *gin.Context) {
	// 	if ctx.Request.FormValue("state") != "random string" {
	// 		fmt.Println("state is not valid")
	// 	}

	// 	token, err := oauth2Conf.Exchange(ctx, ctx.Request.FormValue("code"))
	// 	if err != nil {
	// 		fmt.Printf("could not get token: %s\n", err.Error())
	// 		return
	// 	}

	// 	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	// 	if err != nil {
	// 		fmt.Printf("could not create get request: %s\n", err.Error())
	// 		return
	// 	}
	// 	defer func(Body io.ReadCloser) {
	// 		err := Body.Close()
	// 		if err != nil {
	// 			fmt.Printf("Close response error: %s", err.Error())
	// 			return
	// 		}
	// 	}(resp.Body)
	// 	content, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		fmt.Println("response read error: %s ", err.Error())
	// 		return
	// 	}
	// 	ctx.Data(http.StatusOK, "application/json", content)
	// 	//fmt.Fprintf(ctx.Writer, "Response: %s", content)
	// })

	// log.Fatal(router.Run(":8080"))
}

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var user = User{
	ID:       1,
	Username: "username",
	Password: "password",
}

/*func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid JSON provided.")
		return
	}

	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details.")
		return
	}
	ts, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	saveErr := CreateAuth(user.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
		return
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}*/

/*func CreateToken(userID uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.NewString()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewString()

	var err error

	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	err = os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf")
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}*/

var client *redis.Client

func init() {
	/*dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:16379"
	}*/
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "1234",
		DB:       0,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

/*type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}*/

/*func CreateAuth(userID uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(ctx, td.AccessUUID, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(ctx, td.RefreshUUID, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}*/

type Todo struct {
	UserID uint64 `json:"user_id"`
	Title  string `json:"title"`
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method : %v ", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}

type AccessDetails struct {
	AccessUUID string
	UserID     uint64
}

func FetchAuth(auth *AccessDetails) (uint64, error) {
	accessID, err := client.Get(ctx, auth.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(accessID, 10, 64)
	return userID, nil
}

func CreateTodo(c *gin.Context) {
	var td *Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	tokenAuth, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := FetchAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	td.UserID = userID
	c.JSON(http.StatusCreated, td)
}

func DeleteAuth(uuid string) (int64, error) {
	deleted, err := client.Del(ctx, uuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, err
}

func Logout(c *gin.Context) {
	auth, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	deleted, err := DeleteAuth(auth.AccessUUID)
	if err != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out.")
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

//func Refresh(c *gin.Context) {
//	mapToken := map[string]string{}
//
//	if err := c.ShouldBindJSON(&mapToken); err != nil {
//		c.JSON(http.StatusUnprocessableEntity, err.Error())
//		return
//	}
//
//	refreshToken := mapToken["refresh_token"]
//	//verify the token
//
//	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
//
//	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
//		//Make sure that the token method conform to "SigningMethodHMAC"
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//		}
//		return []byte(os.Getenv("REFRESH_SECRET")), nil
//	})
//
//	//if there is an error, the token must have expired
//	if err != nil {
//		c.JSON(http.StatusUnauthorized, "Refresh token expired")
//		return
//	}
//	//is token valid?
//
//	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
//		c.JSON(http.StatusUnauthorized, err)
//		return
//	}
//	//Since token is valid, get the uuid:
//	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
//
//	if ok && token.Valid {
//		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
//
//		if !ok {
//			c.JSON(http.StatusUnprocessableEntity, err)
//			return
//		}
//
//		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
//
//		if err != nil {
//			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
//			return
//		}
//
//		//Delete the previous Refresh Token
//		deleted, delErr := DeleteAuth(refreshUuid)
//
//		if delErr != nil || deleted == 0 {
//			c.JSON(http.StatusUnauthorized, "unauthorized")
//			return
//		}
//		//Create new pairs of refresh and access tokens
//
//		ts, createErr := CreateToken(userId)
//
//		if createErr != nil {
//			c.JSON(http.StatusForbidden, createErr.Error())
//			return
//		}
//
//		//save the tokens metadata to redis
//		saveErr := CreateAuth(userId, ts)
//
//		if saveErr != nil {
//			c.JSON(http.StatusForbidden, saveErr.Error())
//			return
//		}
//
//		tokens := map[string]string{
//			"access_token":  ts.AccessToken,
//			"refresh_token": ts.RefreshToken,
//		}
//
//		c.JSON(http.StatusCreated, tokens)
//	} else {
//		c.JSON(http.StatusUnauthorized, "refresh expired")
//	}
//}
