package handlers

//// GetClassrooms 获取所有教室
//func GetClassrooms(c *gin.Context) {
//	classrooms, err := services.GetAllClassrooms()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, classrooms)
//}
//
//// GetClassroomByID 根据ID获取教室
//func GetClassroomByID(c *gin.Context) {
//	id := c.Param("id")
//	classroom, err := services.GetClassroomByID(id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Classroom not found"})
//		return
//	}
//	c.JSON(http.StatusOK, classroom)
//}
//
//// AddClassroom 添加教室
//func AddClassroom(c *gin.Context) {
//	var classroom models.Classroom
//	if err := c.ShouldBindJSON(&classroom); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	if err := services.CreateClassroom(&classroom); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, classroom)
//}
//
//// UpdateClassroom 更新教室信息
//func UpdateClassroom(c *gin.Context) {
//	id := c.Param("id")
//	var classroom models.Classroom
//	if err := c.ShouldBindJSON(&classroom); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	if err := services.UpdateClassroom(id, &classroom); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, classroom)
//}
//
//// DeleteClassroom 删除教室
//func DeleteClassroom(c *gin.Context) {
//	id := c.Param("id")
//	if err := services.DeleteClassroom(id); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "Classroom deleted successfully"})
//}
