package services

import (
	"net/http"
	"time"
)

type HeartbeatOperations struct {
	RestOperations
}

// NewHeartbeatOperations creates new instance HeartbeatOperations
// Returns: *HeartbeatOperations
func NewHeartbeatOperations() *HeartbeatOperations {
	hbo := HeartbeatOperations{}
	return &hbo
}

// Heartbeat method are insert timestamp into HTTP result
func (c *HeartbeatOperations) GetHeartbeatOperation() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		c.Heartbeat(res, req)
	}
}

// Heartbeat method are insert timestamp into HTTP result
func (c *HeartbeatOperations) Heartbeat(res http.ResponseWriter, req *http.Request) {
	c.SendResult(res, req, time.Now(), nil)
}
