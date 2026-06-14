package controllers

import (
	"backend/services"
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid emergency ID"})
		return
	}

	// Fetch the emergency request
	emergency, err := tc.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Emergency request not found"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Patient destination coordinates
	patientLat := emergency.Latitude
	patientLng := emergency.Longitude

	// Generate ambulance start position ~2-3km away from patient
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	angle := rng.Float64() * 2 * math.Pi
	// ~0.018-0.027 degrees ≈ 2-3km
	distance := 0.018 + rng.Float64()*0.009
	ambulanceLat := patientLat + distance*math.Cos(angle)
	ambulanceLng := patientLng + distance*math.Sin(angle)

	totalSteps := 15 + rng.Intn(6) // 15-20 steps

	// Send initial dispatch message
	initialMsg := map[string]interface{}{
		"type":    "status_update",
		"status":  "Dispatched",
		"message": "Ambulans telah dikirim ke lokasi Anda",
	}
	if err := conn.WriteJSON(initialMsg); err != nil {
		log.Printf("WebSocket write error: %v", err)
		return
	}

	// Update status to Dispatched in DB
	_ = tc.service.UpdateStatus(uint(id), "Dispatched")

	// Simulate ambulance movement
	for step := 1; step <= totalSteps; step++ {
		time.Sleep(2 * time.Second)

		// Linear interpolation toward patient location
		progress := float64(step) / float64(totalSteps)
		currentLat := ambulanceLat + (patientLat-ambulanceLat)*progress
		currentLng := ambulanceLng + (patientLng-ambulanceLng)*progress

		// Estimate remaining minutes based on remaining steps (each step ≈ 2 seconds simulated as ~1 min real driving)
		remainingSteps := totalSteps - step
		estimatedMinutes := remainingSteps

		// Update tracking in DB
		_ = tc.service.UpdateTracking(uint(id), currentLat, currentLng, estimatedMinutes, "Dispatched")

		// Send location update to client
		locationMsg := map[string]interface{}{
			"type":              "location_update",
			"ambulance_lat":     currentLat,
			"ambulance_lng":     currentLng,
			"estimated_minutes": estimatedMinutes,
			"status":            "Dispatched",
		}
		if err := conn.WriteJSON(locationMsg); err != nil {
			log.Printf("WebSocket write error (client disconnected?): %v", err)
			return
		}
	}

	// Ambulance has arrived
	_ = tc.service.UpdateTracking(uint(id), patientLat, patientLng, 0, "Arrived")

	arrivedMsg := map[string]interface{}{
		"type":    "status_update",
		"status":  "Arrived",
		"message": "Ambulans telah tiba di lokasi Anda",
	}
	if err := conn.WriteJSON(arrivedMsg); err != nil {
		log.Printf("WebSocket write error: %v", err)
		return
	}

	// Close connection gracefully
	conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Ambulance arrived"))
}
