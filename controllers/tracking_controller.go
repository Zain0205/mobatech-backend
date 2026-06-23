package controllers

import (
	"backend/constants"
	"backend/services"
	"backend/utils"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type TrackingController struct {
	service services.EmergencyService
}

func NewTrackingController(service services.EmergencyService) *TrackingController {
	return &TrackingController{service}
}

func (tc *TrackingController) TrackAmbulance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(utils.NewValidationError("Invalid emergency ID"))
		return
	}

	emergency, err := tc.service.GetByID(uint(id))
	if err != nil {
		c.Error(utils.NewAppError(utils.ErrNotFound, http.StatusNotFound, "Emergency request not found"))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	tc.startSimulation(conn, uint(id), emergency.Latitude, emergency.Longitude)
}

func (tc *TrackingController) startSimulation(conn *websocket.Conn, reqID uint, patientLat, patientLng float64) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	angle := rng.Float64() * 2 * math.Pi
	distance := constants.SimulatedDistance + rng.Float64()*constants.SimulatedDistanceRand

	ambLat := patientLat + distance*math.Cos(angle)
	ambLng := patientLng + distance*math.Sin(angle)
	totalSteps := 15 + rng.Intn(6)

	tc.sendInitialStatus(conn, reqID)
	tc.runMovementLoop(conn, reqID, totalSteps, patientLat, patientLng, ambLat, ambLng)
	tc.finalizeArrival(conn, reqID, patientLat, patientLng)
}

func (tc *TrackingController) sendInitialStatus(conn *websocket.Conn, reqID uint) {
	tc.service.UpdateStatus(reqID, constants.StatusDispatched)
	conn.WriteJSON(map[string]interface{}{
		"type":    "status_update",
		"status":  constants.StatusDispatched,
		"message": constants.MsgAmbulanceDispatched,
	})
}

func (tc *TrackingController) runMovementLoop(conn *websocket.Conn, reqID uint, totalSteps int, pLat, pLng, aLat, aLng float64) {
	for step := 1; step <= totalSteps; step++ {
		time.Sleep(constants.SimulatedTimeSleepSec * time.Second)
		progress := float64(step) / float64(totalSteps)

		curLat := aLat + (pLat-aLat)*progress
		curLng := aLng + (pLng-aLng)*progress
		remaining := totalSteps - step

		tc.service.UpdateTracking(reqID, curLat, curLng, remaining, constants.StatusDispatched)
		err := conn.WriteJSON(map[string]interface{}{
			"type":              "location_update",
			"ambulance_lat":     curLat,
			"ambulance_lng":     curLng,
			"estimated_minutes": remaining,
			"status":            constants.StatusDispatched,
		})

		if err != nil {
			log.Printf("WS write error: %v", err)
			return
		}
	}
}

func (tc *TrackingController) finalizeArrival(conn *websocket.Conn, reqID uint, pLat, pLng float64) {
	tc.service.UpdateTracking(reqID, pLat, pLng, 0, "Arrived")
	conn.WriteJSON(map[string]interface{}{
		"type":    "status_update",
		"status":  "Arrived",
		"message": "Ambulans telah tiba di lokasi Anda",
	})
	conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Ambulance arrived"))
}
