package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	userModel "github.com/CAUSALITY-3/Thanal-GO/models/user"
	"github.com/CAUSALITY-3/Thanal-GO/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user userModel.User
	filter := bson.M{"email": email}

	err := s.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		// handle error
	}
	c.Cookie(&fiber.Cookie{
		Name:     "user",
		Value:    string(userJSON),
		MaxAge:   3600000,
		Path:     "/",
		HTTPOnly: false,
		Secure:   false,
	})
	return c.JSON(fiber.Map{"user": user})
}

func (s *UserService) GetAllUsers(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var results []userModel.User

	cursor, err := s.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	defer cursor.Close(ctx)
	// Iterate through the cursor and decode each document into results slice
	for cursor.Next(ctx) {
		var user userModel.User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}
		results = append(results, user)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()

	result, err := s.UserCollection.InsertOne(ctx, &user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	insertedUser := s.UserCollection.FindOne(ctx, bson.M{"_id": result.InsertedID})
	if err := insertedUser.Decode(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reqBody.Update["updatedAt"] = time.Now()
	filter := bson.M{}
	for k, v := range reqBody.Filter {
		filter[k] = v
	}
	update := bson.M{"$set": reqBody.Update}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result userModel.User
	err := s.UserCollection.FindOneAndUpdate(ctx, &filter, &update, opts).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})

	}

	utils.UpdateUsersCache(result)
	return c.JSON(fiber.Map{"result": result})
}

func (s *UserService) UpdateUserOrder(c *fiber.Ctx) error {
	type UpdateRequest struct {
		Filter map[string]string      `json:"filter"` // Criteria to match the document
		Update map[string]interface{} `json:"update"` // Update content
	}

	var reqBody UpdateRequest
	// Get the request body
	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}
	email := reqBody.Filter["email"]
	userData := utils.GetUserCache(email)
	log.Println(userData)
	return nil
}

func (s *UserService) GetUsersCache(c *fiber.Ctx) error {
	usersCache := utils.SingletonInjector.Get("usersCache").(map[string]*userModel.User)
	return c.JSON(fiber.Map{"usersCache": usersCache})
}
