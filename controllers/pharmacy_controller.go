package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PharmacyController struct {
	service services.PharmacyService
}

func NewPharmacyController(service services.PharmacyService) *PharmacyController {
	return &PharmacyController{service}
}

// ==========================================
// PUBLIC & USER HANDLERS
// ==========================================

func (c *PharmacyController) GetCategories(ctx *gin.Context) {
	cats, err := c.service.GetAllCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cats)
}

func (c *PharmacyController) GetMedicines(ctx *gin.Context) {
	catIDStr := ctx.Query("category_id")
	search := ctx.Query("search")

	var catID uint
	if catIDStr != "" {
		parsed, _ := strconv.Atoi(catIDStr)
		catID = uint(parsed)
	}

	meds, err := c.service.GetAllMedicines(catID, search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, meds)
}

func (c *PharmacyController) GetMedicineDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	med, err := c.service.GetMedicineByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Medicine not found"})
		return
	}
	ctx.JSON(http.StatusOK, med)
}

func (c *PharmacyController) GetMyPrescriptions(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64))

	prescriptions, err := c.service.GetPrescriptionsByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, prescriptions)
}

func (c *PharmacyController) GetPrescriptionDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	p, err := c.service.GetPrescriptionByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Prescription not found"})
		return
	}

	// Verify ownership
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64))
	if p.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	ctx.JSON(http.StatusOK, p)
}

func (c *PharmacyController) CreateOrder(ctx *gin.Context) {
	var req models.PharmacyOrder
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	req.UserID = uint(userIDVal.(float64))

	if err := c.service.CreateOrder(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, req)
}

func (c *PharmacyController) GetMyOrders(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64))

	orders, err := c.service.GetOrdersByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

func (c *PharmacyController) GetOrderDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	order, err := c.service.GetOrderByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Verify ownership
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64))
	if order.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (c *PharmacyController) CancelOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	order, err := c.service.GetOrderByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Verify ownership
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64))
	if order.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if order.Status != "Pending" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only pending orders can be cancelled"})
		return
	}

	if err := c.service.UpdateOrderStatus(uint(id), "Cancelled"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

func (c *PharmacyController) GetCart(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64))

	cart, err := c.service.GetCartByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cart)
}

func (c *PharmacyController) AddToCart(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64))

	var req struct {
		MedicineID uint `json:"medicine_id"`
		Quantity   int  `json:"quantity"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.AddToCart(userID, req.MedicineID, req.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Added to cart successfully"})
}

func (c *PharmacyController) UpdateCartItem(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64))

	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateCartItemQuantity(userID, uint(id), req.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Cart item updated successfully"})
}

func (c *PharmacyController) RemoveFromCart(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64))

	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := c.service.RemoveFromCart(userID, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Removed from cart successfully"})
}

// ==========================================
// ADMIN HANDLERS
// ==========================================

func (c *PharmacyController) AdminCreateCategory(ctx *gin.Context) {
	var cat models.MedicineCategory
	if err := ctx.ShouldBindJSON(&cat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.CreateCategory(&cat); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, cat)
}

func (c *PharmacyController) AdminUpdateCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var cat models.MedicineCategory
	if err := ctx.ShouldBindJSON(&cat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.ID = uint(id)
	
	if err := c.service.UpdateCategory(&cat); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cat)
}

func (c *PharmacyController) AdminDeleteCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := c.service.DeleteCategory(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

func (c *PharmacyController) AdminCreateMedicine(ctx *gin.Context) {
	var med models.Medicine
	if err := ctx.ShouldBindJSON(&med); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.CreateMedicine(&med); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, med)
}

func (c *PharmacyController) AdminUpdateMedicine(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var med models.Medicine
	if err := ctx.ShouldBindJSON(&med); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	med.ID = uint(id)

	if err := c.service.UpdateMedicine(&med); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, med)
}

func (c *PharmacyController) AdminDeleteMedicine(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := c.service.DeleteMedicine(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Medicine deleted"})
}

func (c *PharmacyController) AdminCreatePrescription(ctx *gin.Context) {
	var p models.Prescription
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.CreatePrescription(&p); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, p)
}

func (c *PharmacyController) AdminGetAllPrescriptions(ctx *gin.Context) {
	prescriptions, err := c.service.GetAllPrescriptions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, prescriptions)
}

func (c *PharmacyController) AdminGetAllOrders(ctx *gin.Context) {
	orders, err := c.service.GetAllOrders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

func (c *PharmacyController) AdminUpdateOrderStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateOrderStatus(uint(id), req.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func (c *PharmacyController) AdminUpdateOrderPayment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req struct {
		PaymentStatus string `json:"payment_status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateOrderPayment(uint(id), req.PaymentStatus); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Payment status updated successfully"})
}
