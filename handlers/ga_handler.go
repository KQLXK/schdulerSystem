package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"schedule/algorithm"
	"schedule/commen/result"
	"schedule/dto"
)

func GAHandler(c *gin.Context) {
	var req dto.ScheduleGARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}
	log.Println("schedule/ga: get req success, req:", req)
	resp := algorithm.GASchedule(&req.AdditionalRules, req.ScheduleIDs)
	result.Sucess(c, resp)
}
