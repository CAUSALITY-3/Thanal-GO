package initializers

import (
	database "github.com/CAUSALITY-3/Thanal-GO/models/DB"
	services "github.com/CAUSALITY-3/Thanal-GO/service/user"
	"github.com/CAUSALITY-3/Thanal-GO/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func InjectServices() error {
	dbClient := database.GetDBClient()
	userService := services.NewUserService(getCollection(dbClient, "users"))
	utils.SingletonInjector.Bind(userService, "userService")
	return nil
}

func getCollection(client *mongo.Client, collection string) *mongo.Collection {
	return client.Database("thanal").Collection(collection)
}

// func injectServices() error {
//     inj := injector.NewInjector()

//     productFeatureService := &services.ProductFeatureService{
//         ProductFeatures: models.ProductFeature{},
//     }
//     if err := inj.Bind(productFeatureService, "productFeatureService"); err != nil {
//         return err
//     }

//     productService := &services.ProductService{
//         Product:               models.Product{},
//         ProductMainList:       models.ProductMainList{},
//         ProductFeatureService: productFeatureService,
//     }
//     if err := inj.Bind(productService, "productService"); err != nil {
//         return err
//     }

//     userService := &services.UserService{
//         User: models.User{},
//     }
//     if err := inj.Bind(userService, "userService"); err != nil {
//         return err
//     }

//     authenticationService := &services.AuthenticationService{
//         UserService:        userService,
//         GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
//         GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
//         GoogleCallbackURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
//     }
//     if err := inj.Bind(authenticationService, "authenticationService"); err != nil {
//         return err
//     }

//     imageService := &services.ImageService{}
//     if err := inj.Bind(imageService, "imageService"); err != nil {
//         return err
//     }

//     // Mock call to get all users, assuming userService has a method GetAllUsers
//     allUsers, err := userService.GetAllUsers(context.Background())
//     if err != nil {
//         return err
//     }
//     usersCache := make(map[string]*models.User)
//     for _, user := range allUsers {
//         usersCache[user.Email] = user
//     }
//     if err := inj.Bind(usersCache, "usersCache"); err != nil {
//         return err
//     }

//     razorpayConfig := map[string]string{
//         "key_id":     os.Getenv("RAZORPAY_KEY"),
//         "key_secret": os.Getenv("RAZORPAY_SECRET"),
//     }
//     paymentService := &services.PaymentService{
//         RazorpayConfig: razorpayConfig,
//     }
//     if err := inj.Bind(paymentService, "paymentService"); err != nil {
//         return err
//     }

//     uploadService := &services.UploadService{}
//     if err := inj.Bind(uploadService, "uploadService"); err != nil {
//         return err
//     }

//     orderService := &services.OrderService{
//         Order:         models.Order{},
//         ProductService: productService,
//         UserService:   userService,
//     }
//     if err := inj.Bind(orderService, "orderService"); err != nil {
//         return err
//     }

//     // Assuming loadCache is a function that performs necessary cache operations
//     if err := loadCache(); err != nil {
//         return err
//     }

//     log.Println("Services injected successfully")
//     return nil
// }

// func main() {
//     if err := injectServices(); err != nil {
//         log.Fatalf("Failed to inject services: %v", err)
//     }
//     // Continue with server initialization or other operations...
// }
