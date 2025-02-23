package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type SloganController struct {
	ssr service.SloganService //依赖与接口
}

func NewSloganController(ssr service.SloganService) *SloganController {
	return &SloganController{
		ssr: ssr,
	}
}

func (uc *SloganController) GetSlogan(c *gin.Context) {
	device := c.Param("device_num")
	if device == "" {
		c.JSON(http.StatusBadRequest, response.Response{Error: "设备号不能为空"})
		return
	}

	err := uc.ssr.GetSlogan(device)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
	}
	c.JSON(http.StatusOK, "获取激励语成功")
}

func (uc *SloganController) ChangeSlogan(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Response{Error: "无该用户"})
		return
	}
	var newSlogan request.Slogan
	if err := c.ShouldBindJSON(&newSlogan); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "识别标语失败"})
		return
	}
	//显式转换
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}

	err = uc.ssr.ChangeSlogan(uint(idInt), newSlogan)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Message: "激励语更新成功"})
}
