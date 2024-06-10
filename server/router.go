package server

import (
	configs "beli_mang/cfg"
	"beli_mang/controller"
	"beli_mang/middleware"
	"beli_mang/repo"
	"beli_mang/svc"
	"context"

	"log"

	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (s *Server) RegisterRoute(config configs.Config) {
	mainRoute := s.app.Group("")

	registerAdminRoute(mainRoute, s.dbPool)
	registerUserRoute(mainRoute, s.dbPool)
	registerMerchantRoute(mainRoute, s.dbPool)
	registerPurchaseRoute(mainRoute, s.dbPool)
	registerImageRoute(mainRoute, config)
}

func registerAdminRoute(r fiber.Router, db *pgxpool.Pool) {
	ctr := controller.NewUserController(svc.NewUserSvc(repo.NewUserRepo(db)))

	adminGroup := r.Group("/admin")
	newRoute(adminGroup, "POST", "/register", ctr.AdminRegister)
	newRoute(adminGroup, "POST", "/login", ctr.AdminLogin)
}

func registerUserRoute(r fiber.Router, db *pgxpool.Pool) {
	ctr := controller.NewUserController(svc.NewUserSvc(repo.NewUserRepo(db)))

	user := r.Group("/users")
	newRoute(user, "POST", "/register", ctr.UserRegister)
	newRoute(user, "POST", "/login", ctr.UserLogin)
}

func registerMerchantRoute(r fiber.Router, db *pgxpool.Pool) {
	ctr := controller.NewMerchantController(svc.NewMerchantSvc(repo.NewMerchantRepo(db)))

	merchantGroup := r.Group("/admin/merchants")
	newRouteWithAdminAuth(merchantGroup, "POST", "", ctr.RegisterMerchant)
	newRouteWithAdminAuth(merchantGroup, "GET", "", ctr.GetMerchant)
	newRouteWithAdminAuth(merchantGroup, "POST", "/:merchantId/items", ctr.RegisterItem)
	newRouteWithAdminAuth(merchantGroup, "GET", "/:merchantId/items", ctr.GetItem)
}

func registerPurchaseRoute(r fiber.Router, db *pgxpool.Pool) {
	ctr := controller.NewPurchaseController(svc.NewPurchaseSvc(repo.NewPurchaseRepo(db)))

	merchantsGroup := r.Group("/merchants")
	newRouteWithUserAuth(merchantsGroup, "GET", "/nearby/:latlong", ctr.GetNearbyMerchant)

	usersGroup := r.Group("/users")
	newRouteWithUserAuth(usersGroup, "POST", "/estimate", ctr.EstimateOrder)
	newRouteWithUserAuth(usersGroup, "POST", "/orders", ctr.PlaceOrder)
	newRouteWithUserAuth(usersGroup, "GET", "/orders", ctr.GetOrder)
}

func registerImageRoute(r fiber.Router, config configs.Config) {
	bucket := config.AWS_S3_BUCKET_NAME

	// Load AWS configuration
	cfg, err := awsCfg.LoadDefaultConfig(
		context.Background(),
		awsCfg.WithRegion(config.AWS_REGION),
		awsCfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.AWS_ACCESS_KEY_ID,
				config.AWS_SECRET_ACCESS_KEY,
				"",
			),
		),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	// Initialize the image service and controller
	imageService := svc.NewImageSvc(cfg, bucket)
	imageController := controller.NewImageController(imageService)

	// Register the route with authentication middleware
	newRouteWithAdminAuth(r, "POST", "/image", imageController.UploadImage)
}

func newRoute(router fiber.Router, method, path string, handler fiber.Handler) {
	router.Add(method, path, handler)
}

func newRouteWithAuth(router fiber.Router, method, path string, handler fiber.Handler) {
	router.Add(method, path, middleware.Auth, handler)
}

func newRouteWithAdminAuth(router fiber.Router, method, path string, handler fiber.Handler) {
	router.Add(method, path, middleware.AdminAuth, handler)
}

func newRouteWithUserAuth(router fiber.Router, method, path string, handler fiber.Handler) {
	router.Add(method, path, middleware.UserAuth, handler)
}
