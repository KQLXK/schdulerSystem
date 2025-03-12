package handlers

//// GetClasses 获取所有班级
//func GetClasses(c *gin.Context) {
//	classes, err := services.GetAllClasses()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, classes)
//}
//
//// GetClassByID 根据ID获取班级
//func GetClassByID(c *gin.Context) {
//	id := c.Param("id")
//	class, err := services.GetClassByID(id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
//		return
//	}
//	c.JSON(http.StatusOK, class)
//}
//
//// AddClass 添加班级
//func AddClass(c *gin.Context) {
//	var class models.Class
//	if err := c.ShouldBindJSON(&class); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	if err := services.CreateClass(&class); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, class)
//}
//
//// UpdateClass 更新班级信息
//func UpdateClass(c *gin.Context) {
//	id := c.Param("id")
//	var class models.Class
//	if err := c.ShouldBindJSON(&class); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	if err := services.UpdateClass(id, &class); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, class)
//}
//
//// DeleteClass 删除班级
//func DeleteClass(c *gin.Context) {
//	id := c.Param("id")
//	if err := services.DeleteClass(id); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "Class deleted successfully"})
//}
