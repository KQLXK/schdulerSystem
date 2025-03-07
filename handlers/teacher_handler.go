package handlers

//// GetTeachers 获取所有教师
//func GetTeachers(c *gin.Context) {
//	teachers, err := services.GetAllTeachers()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, teachers)
//}
//
//// GetTeacherByID 根据ID获取教师
//func GetTeacherByID(c *gin.Context) {
//	id := c.Param("id")
//	teacher, err := services.GetTeacherByID(id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
//		return
//	}
//	c.JSON(http.StatusOK, teacher)
//}
//
//// AddTeacher 添加教师
//func AddTeacher(c *gin.Context) {
//	var teacher models.Teacher
//	if err := c.ShouldBindJSON(&teacher); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	if err := services.CreateTeacher(&teacher); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, teacher)
//}
//
//// UpdateTeacher 更新教师信息
//func UpdateTeacher(c *gin.Context) {
//	id := c.Param("id")
//	var teacher models.Teacher
//	if err := c.ShouldBindJSON(&teacher); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	if err := services.UpdateTeacher(id, &teacher); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, teacher)
//}
//
//// DeleteTeacher 删除教师
//func DeleteTeacher(c *gin.Context) {
//	id := c.Param("id")
//	if err := services.DeleteTeacher(id); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "Teacher deleted successfully"})
//}
