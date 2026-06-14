package routes

import (
	"backend/controllers"
	"backend/middleware"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPharmacyRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repositories.NewPharmacyRepository(db)
	service := services.NewPharmacyService(repo)
	ctrl := controllers.NewPharmacyController(service)

	// Public routes (no auth needed for browsing)
	public := r.Group("/api/pharmacy")
	{
		public.GET("/categories", ctrl.GetCategories)
		public.GET("/medicines", ctrl.GetMedicines)
		public.GET("/medicines/:id", ctrl.GetMedicineDetail)
	}

	// User routes (auth required)
	user := r.Group("/api/pharmacy")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/prescriptions", ctrl.GetMyPrescriptions)
		user.GET("/prescriptions/:id", ctrl.GetPrescriptionDetail)
		user.POST("/orders", ctrl.CreateOrder)
		user.GET("/orders", ctrl.GetMyOrders)
		user.GET("/orders/:id", ctrl.GetOrderDetail)
		user.PUT("/orders/:id/cancel", ctrl.CancelOrder)

		user.GET("/cart", ctrl.GetCart)
		user.POST("/cart", ctrl.AddToCart)
		user.PUT("/cart/:id", ctrl.UpdateCartItem)
		user.DELETE("/cart/:id", ctrl.RemoveFromCart)
	}

	// Admin routes
	admin := r.Group("/api/admin/pharmacy")
	{
		admin.POST("/categories", ctrl.AdminCreateCategory)
		admin.PUT("/categories/:id", ctrl.AdminUpdateCategory)
		admin.DELETE("/categories/:id", ctrl.AdminDeleteCategory)
		admin.POST("/medicines", ctrl.AdminCreateMedicine)
		admin.PUT("/medicines/:id", ctrl.AdminUpdateMedicine)
		admin.DELETE("/medicines/:id", ctrl.AdminDeleteMedicine)
		admin.POST("/prescriptions", ctrl.AdminCreatePrescription)
		admin.GET("/prescriptions", ctrl.AdminGetAllPrescriptions)
		admin.GET("/orders", ctrl.AdminGetAllOrders)
		admin.PUT("/orders/:id/status", ctrl.AdminUpdateOrderStatus)
		admin.PUT("/orders/:id/payment", ctrl.AdminUpdateOrderPayment)
	}
}
