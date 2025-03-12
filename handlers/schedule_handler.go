package handlers

//// GetSchedules 获取所有排课结果
//func GetSchedules(c *gin.Context) {
//	schedules, err := services.GetAllSchedules()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, schedules)
//}
//
//// GetScheduleByID 根据ID获取排课结果
//func GetScheduleByID(c *gin.Context) {
//	id := c.Param("id")
//	schedule, err := services.GetScheduleByID(id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
//		return
//	}
//	c.JSON(http.StatusOK, schedule)
//}
//
//// AddSchedule 添加排课结果
//func AddSchedule(c *gin.Context) {
//	var schedule models.Schedule
//	if err := c.ShouldBindJSON(&schedule); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	if err := services.CreateSchedule(&schedule); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, schedule)
//}
//
//// UpdateSchedule 更新排课结果
//func UpdateSchedule(c *gin.Context) {
//	id := c.Param("id")
//	var schedule models.Schedule
//	if err := c.ShouldBindJSON(&schedule); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	if err := services.UpdateSchedule(id, &schedule); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, schedule)
//}
//
//// DeleteSchedule 删除排课结果
//func DeleteSchedule(c *gin.Context) {
//	id := c.Param("id")
//	if err := services.DeleteSchedule(id); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
//}
