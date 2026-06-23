package controllers

import (
	"backend/services"
	"backend/utils"
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

func (c *PharmacyController) GetCategories(ctx *gin.Context) {
	cats, err := c.service.GetAllCategories()
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", cats))
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
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", meds))
}

func (c *PharmacyController) GetMedicineDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	med, err := c.service.GetMedicineByID(uint(id))
	if err != nil {
		ctx.Error(utils.NewAppError(utils.ErrNotFound, http.StatusNotFound, "Medicine not found"))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", med))
}
