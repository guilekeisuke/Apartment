package controller

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"api/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/cyruzin/golang-tmdb"
)

type jwtCustomClaims struct {
	UID  int    `json:"uid"`
	Name string `json:"name"`
	jwt.StandardClaims
}

var signingKey = []byte("secret")

var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

func getTmdbInfo(c echo.Context) error {
	tmdbClient, err := tmdb.Init(os.GetEnv("YOUR_APIKEY"))
	if err != nil {
		fmt.Println(err)
	}

	// OPTIONAL: Setting a custom config for the http.Client.
	// The default timeout is 10 seconds. Here you can set other
	// options like Timeout and Transport.
	customClient := http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns: 10,
			IdleConnTimeout: 15 * time.Second,
		}
	}

	tmdbClient.SetClientConfig(customClient)

	// OPTIONAL: Enable this option if you're going to use endpoints
	// that needs session id.
	// 
	// You can read more about how this works:
	// https://developers.themoviedb.org/3/authentication/how-do-i-generate-a-session-id
	tmdbClient.SetSessionID(os.GetEnv("YOUR_SESSION_ID"))

	// OPTIONAL (Recommended): Enabling auto retry functionality.
	// This option will retry if the previous request fail.
	tmdbClient.SetClientAutoRetry()

	movie, err := tmdbClient.GetMovieDetails(297802, nil)
	if err != nil {
	fmt.Println(err)
	}

	fmt.Println(movie.Title)
	return c.JSON(http.StatusOK, movie.Title)
}

func Signup(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	if user.UserName == "" || user.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid name or password",
		}
	}

	if u := model.FindUser(&model.User{UserName: user.UserName}); u.ID != 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "name already exists",
		}
	}

	hash, err := passwordHash(user.Password)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "cannot created hash password",
		}
	}
	inputPassword := user.Password
	user.Password = hash
	model.CreateUser(user)
	user.Password = inputPassword

	return c.JSON(http.StatusCreated, user)
}

func Login(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	user := model.FindUser(&model.User{UserName: u.UserName})

	err := passwordVerify(user.Password, u.Password)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "invalid name or password",
		}
	}

	claims := &jwtCustomClaims{
		user.ID,
		user.UserName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(signingKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func userIDFromToken(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	uid := claims.UID
	return uid
}

func passwordHash(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

// パスワードがハッシュにマッチするかどうかを調べる
func passwordVerify(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
