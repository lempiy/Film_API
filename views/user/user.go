package user

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"time"
	"github.com/lempiy/echo_api/models"
	"fmt"
	"github.com/lempiy/echo_api/utils"
	"github.com/lempiy/echo_api/types"
	"github.com/lempiy/echo_api/utils/validator"
)

type person struct {
	Login     string `json:"login" form:"login" query:"login"`
	Password string `json:"password" form:"password" query:"password"`
}

//Login function used to get jwt token by name and password,
//requests login and password in JSON body.
func Login(c echo.Context) error {
	u := new(person)
	if err := c.Bind(u); err != nil {
		return err
	}
	usr, err := models.User.ReadByLogin(u.Login)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong",
		})
	}
	encryptPass := utils.EncryptPassword(u.Password)
	if u.Login == usr.Login && encryptPass == usr.Password {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = usr.ID
		claims["name"] = u.Login
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
	return c.JSON(http.StatusForbidden, map[string]string{
		"error": "Wrong password or username.",
	})
}

func Register(c echo.Context) error {
	u := new(types.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if isOK, message := validator.ValidateUserData(u); !isOK {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": message,
		})
	}
	err := models.User.Create(u)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong",
		})
	}
	return c.JSON(http.StatusOK, map[string]bool{
		"success": true,
	})
}

func Test(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return c.JSON(http.StatusOK, claims)
}
