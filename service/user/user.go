package services

import (
	"context"
	"log"
	"time"

	database "github.com/CAUSALITY-3/Thanal-GO/models/DB"
	userModel "github.com/CAUSALITY-3/Thanal-GO/models/user"
	"github.com/CAUSALITY-3/Thanal-GO/utils"
	"github.com/gofiber/fiber/v2"
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

func (s *UserService) FindUserByEmail(c *fiber.Ctx) error {
	email := c.Query("email")

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	var user userModel.User
	filter := bson.M{"email": email}

	user, err := database.MongoFindOne(s.UserCollection, user, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// userJSON, err := json.Marshal(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// c.Cookie(&fiber.Cookie{
	// 	Name:     "user",
	// 	Value:    string(userJSON),
	// 	MaxAge:   3600000,
	// 	Path:     "/",
	// 	HTTPOnly: false,
	// 	Secure:   false,
	// })
	return c.JSON(fiber.Map{"user": user})
}

func (s *UserService) GetAllUsers(c *fiber.Ctx) error {

	results, err := database.MongoFindAll(s.UserCollection, userModel.User{})
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.JSON(results)
}

func (s *UserService) CacheAllUsers() bool {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer cursor.Close(ctx)

	usersCache := make(map[string]*userModel.User)
	// Iterate through the cursor and decode each document into results slice
	for cursor.Next(ctx) {
		var user userModel.User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
			return false
		}
		usersCache[user.Email] = &user
	}
	utils.SingletonInjector.Bind(usersCache, "usersCache")
	return true
}

func (s *UserService) CreateUser(c *fiber.Ctx) error {
	var user userModel.User
	// Get the request body
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	// Validate the request body
	if err := utils.ValidateStruct(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := database.MongoCreate(s.UserCollection, user)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	utils.UpdateUsersCache(user)
	return c.JSON(user)
}

func (s *UserService) UpsertUser(c *fiber.Ctx) error {

	type UpdateRequest struct {
		Filter map[string]string      `json:"filter"` // Criteria to match the document
		Update map[string]interface{} `json:"update"` // Update content
	}

	var reqBody UpdateRequest
	// Get the request body
	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}

	reqBody.Update["updatedAt"] = time.Now()
	filter := bson.M{}
	for k, v := range reqBody.Filter {
		filter[k] = v
	}
	update := bson.M{"$set": reqBody.Update}

	result, err := database.MongoFindOneAndUpdate(s.UserCollection, userModel.User{}, filter, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	utils.UpdateUsersCache(result)
	return c.JSON(fiber.Map{"result": result})
}

func (s *UserService) UpdateUserOrder(c *fiber.Ctx) error {
	type UpdateType struct {
		OrderItems map[string]any `json:"orderItems"` //TODO update the type in future
		OrderId    string         `json:"orderId"`
	}
	type UpdateRequest struct {
		Filter map[string]string `json:"filter"` // Criteria to match the document
		Update UpdateType        `json:"update"` // Update content
	}

	var reqBody UpdateRequest
	// Get the request body
	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}
	email := reqBody.Filter["email"]
	userData := utils.GetUserCache(email)
	bag := utils.Filter(userData.Bag, func(id string) bool {
		return reqBody.Update.OrderItems[id] != nil
	})
	orderAlreadyExists := utils.Includes(userData.Orders, func(id string) bool {
		return id == reqBody.Update.OrderId
	})

	if orderAlreadyExists {
		return c.Status(fiber.StatusBadRequest).JSON(userData)
	}

	filter := bson.M{"email": reqBody.Filter["email"]}
	update := bson.M{"$push": bson.M{"orders": reqBody.Update.OrderId}, "$set": bson.M{"bag": bag, "updatedAt": time.Now()}}
	updatedUser, err := database.MongoFindOneAndUpdate(s.UserCollection, userModel.User{}, filter, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	c.JSON(updatedUser)
	return nil
}

func (s *UserService) AddToBag(c *fiber.Ctx) error {

	type UpdateType struct {
		ProductId string `json:"productId"`
	}

	type FilterType struct {
		Email string `json:"email"`
	}
	type UpdateRequest struct {
		Filter FilterType `json:"filter"` // Criteria to match the document
		Update UpdateType `json:"update"` // Update content
	}
	// var reqBody UpdateRequest
	// Get the request body

	reqBody, err := utils.ParseBody[UpdateRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	log.Println("Adding to bag:", reqBody.Update.ProductId)
	email := reqBody.Filter.Email

	userCache := utils.GetUserCache(email)

	if utils.Includes(userCache.Bag, func(id string) bool {
		return id == reqBody.Update.ProductId
	}) {
		return c.Status(fiber.StatusConflict).JSON(userCache)
	}

	filter := bson.M{"email": email}
	update := bson.M{"$push": bson.M{"bag": reqBody.Update.ProductId}, "$set": bson.M{"updatedAt": time.Now()}}

	updatedUser, err := database.MongoFindOneAndUpdate(s.UserCollection, userModel.User{}, filter, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	utils.UpdateUsersCache(updatedUser)
	utils.CookieUpdate(c, updatedUser)
	c.JSON(updatedUser)
	return nil
}

func (s *UserService) RemoveFromBag(c *fiber.Ctx) error {

	type UpdateType struct {
		ProductId string `json:"productId"`
	}

	type FilterType struct {
		Email string `json:"email"`
	}
	type UpdateRequest struct {
		Filter FilterType `json:"filter"` // Criteria to match the document
		Update UpdateType `json:"update"` // Update content
	}
	reqBody, err := utils.ParseBody[UpdateRequest](c)
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	email := reqBody.Filter.Email
	userCache := utils.GetUserCache(email)

	if !utils.Includes(userCache.Bag, func(id string) bool {
		return id == reqBody.Update.ProductId
	}) {
		return c.Status(fiber.StatusAlreadyReported).JSON(userCache)
	}

	filter := bson.M{"email": email}
	update := bson.M{"$pull": bson.M{"bag": reqBody.Update.ProductId}, "$set": bson.M{"updatedAt": time.Now()}}

	updatedUser, err := database.MongoFindOneAndUpdate(s.UserCollection, userModel.User{}, filter, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	utils.UpdateUsersCache(updatedUser)
	utils.CookieUpdate(c, updatedUser)
	c.JSON(updatedUser)
	return nil
}

func (s *UserService) FavoriteItem(c *fiber.Ctx) error {

	type UpdateType struct {
		ProductId string `json:"productId"`
	}

	type FilterType struct {
		Email string `json:"email"`
	}
	type UpdateRequest struct {
		Filter FilterType `json:"filter"` // Criteria to match the document
		Update UpdateType `json:"update"` // Update content
	}
	reqBody, err := utils.ParseBody[UpdateRequest](c)
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	email := reqBody.Filter.Email
	userCache := utils.GetUserCache(email)

	if utils.Includes(userCache.Wishlists, func(id string) bool {
		return id == reqBody.Update.ProductId
	}) {
		return c.Status(fiber.StatusConflict).JSON(userCache)
	}

	filter := bson.M{"email": email}
	update := bson.M{"$push": bson.M{"wishlists": reqBody.Update.ProductId}, "$set": bson.M{"updatedAt": time.Now()}}

	updatedUser, err := database.MongoFindOneAndUpdate(s.UserCollection, userModel.User{}, filter, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	utils.UpdateUsersCache(updatedUser)
	utils.CookieUpdate(c, updatedUser)
	c.JSON(updatedUser)
	return nil
}

func (s *UserService) UnfavoriteItem(c *fiber.Ctx) error {

	type UpdateType struct {
		ProductId string `json:"productId"`
	}

	type FilterType struct {
		Email string `json:"email"`
	}
	type UpdateRequest struct {
		Filter FilterType `json:"filter"` // Criteria to match the document
		Update UpdateType `json:"update"` // Update content
	}
	reqBody, err := utils.ParseBody[UpdateRequest](c)
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	email := reqBody.Filter.Email
	userCache := utils.GetUserCache(email)

	if !utils.Includes(userCache.Wishlists, func(id string) bool {
		return id == reqBody.Update.ProductId
	}) {
		return c.Status(fiber.StatusAlreadyReported).JSON(userCache)
	}

	filter := bson.M{"email": email}
	update := bson.M{"$pull": bson.M{"wishlists": reqBody.Update.ProductId}, "$set": bson.M{"updatedAt": time.Now()}}

	updatedUser, err := database.MongoFindOneAndUpdate(s.UserCollection, userModel.User{}, filter, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	utils.UpdateUsersCache(updatedUser)
	utils.CookieUpdate(c, updatedUser)
	c.JSON(updatedUser)
	return nil
}

func (s *UserService) GetUsersCache(c *fiber.Ctx) error {
	usersCache := utils.SingletonInjector.Get("usersCache").(map[string]*userModel.User)
	return c.JSON(fiber.Map{"usersCache": usersCache})
}
