package utils

import (
	"log"

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

func GetUserCache(email string) *userModel.User {
	usersCache := SingletonInjector.Get("usersCache").(map[string]*userModel.User)
	return usersCache[email]
}

func Filter[T any](slice []T, condition func(T) bool) []T {
	var result []T = []T{}
	for _, item := range slice {
		if condition(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](slice []T, transform func(T) R) []R {
	var result []R = make([]R, len(slice))
	for index, item := range slice {
		result[index] = transform(item)
	}
	return result
}

func Find[T any](slice []T, condition func(T) bool) *T {
	for _, item := range slice {
		if condition(item) {
			return &item
		}
	}
	return nil
}

func Includes[T any](slice []T, condition func(T) bool) bool {
	log.Println("Slice", slice)
	for _, item := range slice {
		if condition(item) {
			return true
		}
	}
	return false
}
