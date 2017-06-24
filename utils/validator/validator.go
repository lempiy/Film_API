package validator

import (
	"github.com/lempiy/echo_api/types"
	"fmt"
	"regexp"
)

func ValidateUserData(u * types.User) (bool, string) {
	phonePattern := regexp.MustCompile(`[0-9\s\+\-\(\)]{5,20}`)
	if len(u.Username) > 30 {
		return false, fmt.Sprintf("Invalid username - %s", u.Username)
	}
	if u.Password == "" || len(u.Password) < 4 {
		return false, fmt.Sprintf("Password should be longer then 4 symbols")
	}
	if u.Login == "" || len(u.Login) < 3 {
		return false, fmt.Sprintf("Login should be longer then 3 symbols")
	}
	if !phonePattern.MatchString(u.Telephone) {
		return false, fmt.Sprintf("Incorrect phone number")
	}
	if u.Age < 16 || u.Age > 120 {
		return false, fmt.Sprintf("Incorrect user age")
	}
	return true, ""
}

func ValidateGenresQuery(genresQuery string) (bool, string) {
	queryParamsPattern := regexp.MustCompile(`^\d+(,\d+)*$`)
	if !queryParamsPattern.MatchString(genresQuery) {
		return false, fmt.Sprintf("Incorrect genre query. Genres should be seperated by commas.")
	}
	return true, ""
}
