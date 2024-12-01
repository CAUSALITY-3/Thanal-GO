package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	userModel "github.com/CAUSALITY-3/Thanal-GO/models/user"
	"github.com/CAUSALITY-3/Thanal-GO/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	UserCollection *mongo.Collection
}

func NewUserService(userCollection *mongo.Collection) *UserService {
	return &UserService{
		UserCollection: userCollection,
	}
}

func (s *UserService) FindUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		email = "abinbabu003@gmail.com"
	}
	log.Println("tttttttttttttttttt", email)
	var body map[string]interface{}
	log.Println("qqqqqqqqqqqqqqqqqqqqqqq")
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("wwwwwwwwwwwwwwwww")
		body = nil
	}
	log.Println("eeeeeeeeeeeee")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user []userModel.User
	filter := bson.M{"email": email}

	err := s.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		// handle error
	}
	c.SetCookie("user", string(userJSON), 3600000, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"user": user, "body": body})
}

func (s *UserService) CreateUser(c *gin.Context) {

	var user userModel.User
	if err := utils.GetReqBody(c, &user); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()

	result, err := s.UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"result": result, "user": user})
}

// func (s *UserService) UpsertUser(ctx context.Context, data bson.M) (*UserData, error) {
// 	email, ok := data["email"].(string)
// 	if !ok || email == "" {
// 		return nil, errors.New("user email not provided")
// 	}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"lastLoggedIn": time.Now(),
// 			"updatedAt":    time.Now(),
// 			"name":         data["name"],
// 			"email":        email,
// 			"profilePic":   data["picture"],
// 		},
// 	}

// 	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
// 	var user UserData
// 	err := s.UserCollection.FindOneAndUpdate(ctx, bson.M{"email": email}, update, opts).Decode(&user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (s *UserService) UpdateUserByQuery(ctx context.Context, query, data bson.M) (*UserData, error) {
// 	update := bson.M{
// 		"$set": bson.M{
// 			"updatedAt": time.Now(),
// 		},
// 	}
// 	for key, value := range data {
// 		update["$set"].(bson.M)[key] = value
// 	}

// 	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
// 	var user UserData
// 	err := s.UserCollection.FindOneAndUpdate(ctx, query, update, opts).Decode(&user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (s *UserService) UpdateUserOrder(ctx context.Context, query bson.M, orderId string, orderItems map[string]bool) (*UserData, error) {

// 	UsersCache, ok := utils.SingletonInjector.Get("UsersCache").(map[string]*UserData)
// 	if !ok {
// 		return nil, errors.New("UsersCache not found")
// 	}

// 	userdata, exists := UsersCache[query["email"].(string)]
// 	if !exists {
// 		var result *mongo.SingleResult
// 		err := s.UserCollection.FindOne(ctx, query).Decode(&result)
// 		if err != nil {
// 			return nil, err
// 		}
// 		var userdata UserData
// 		err = result.Decode(&userdata)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	bag := []string{}
// 	for _, item := range userdata.Bag {
// 		if !orderItems[item] {
// 			bag = append(bag, item)
// 		}
// 	}

// 	if contains(userdata.Orders, orderId) {
// 		return userdata, nil
// 	}

// 	update := bson.M{
// 		"$push": bson.M{"orders": orderId},
// 		"$set":  bson.M{"bag": bag, "updatedAt": time.Now()},
// 	}

// 	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
// 	var user UserData
// 	err := s.UserCollection.FindOneAndUpdate(ctx, query, update, opts).Decode(&user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	userdata[user.Email] = &user
// 	return &user, nil
// }

// func (s *UserService) AddToBag(ctx context.Context, query bson.M, productId string) (*UserData, error) {
// 	userdata, exists := s.UsersCache[query["email"].(string)]
// 	if !exists {
// 		return nil, errors.New("user not found in cache")
// 	}

// 	if contains(userdata.Bag, productId) {
// 		return userdata, nil
// 	}

// 	update := bson.M{
// 		"$push": bson.M{"bag": productId},
// 		"$set":  bson.M{"updatedAt": time.Now()},
// 	}

// 	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
// 	var user UserData
// 	err := s.UserCollection.FindOneAndUpdate(ctx, query, update, opts).Decode(&user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	s.UsersCache[user.Email] = &user
// 	return &user, nil
// }

// // Helper function to check if a slice contains a value
// func contains(slice []string, value string) bool {
// 	for _, item := range slice {
// 		if item == value {
// 			return true
// 		}
// 	}
// 	return false
// }

// Additional functions for removeFromBag, favoriteItem, unFavoriteItem, etc., would follow a similar pattern.
