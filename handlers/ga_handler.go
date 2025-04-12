package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"schedule/algorithm"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/services/table"
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

func ManualScheduleHandler(c *gin.Context) {
	var req dto.ManualScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}
	log.Println("/table/create: get req success, req:", req)
	err := table.ManualSchedule(&req)
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, nil)
}

func UpdateTableHandler(c *gin.Context) {
	var req dto.ScheduleResultUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}
	log.Println("/table/update: get req success, req:", req)
	err := table.UpdateTable(&req)
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, nil)
}
