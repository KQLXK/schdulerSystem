package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"schedule/commen/result"
	"schedule/dto"
)

func UserLoginHandler(c *gin.Context) {
	var u dto.UserLoginRequest
	err := c.ShouldBind(&u)
	if err != nil {
		log.Println("get userinfo failed, err: ", err)
		result.Error(c, result.GetReqErrStatus)
		return
	}

	//data, err := user.UserLogin(u)
	if err != nil {
	} else {
		result.Sucess(c, nil)
		return
	}
}
