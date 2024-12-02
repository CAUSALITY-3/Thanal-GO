package utils

import (
	userModel "github.com/CAUSALITY-3/Thanal-GO/models/user"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func UpdateUsersCache(user userModel.User) bool {
	usersCache := SingletonInjector.Get("usersCache").(map[string]*userModel.User)
	usersCache[user.Email] = &user
	return SingletonInjector.Update(usersCache, "usersCache")
}

func GetUserCache(email string) userModel.User {
	usersCache := SingletonInjector.Get("usersCache").(map[string]*userModel.User)
	return *usersCache[email]
}
