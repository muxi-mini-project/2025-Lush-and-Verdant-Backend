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

// GetSlogan 获取激励语
// @Summary 获取激励语
// @Description 根据设备号获取激励语
// @Tags 标语
// @Accept json
// @Produce json
// @Param device_num path string true "设备号"
// @Success 200 {object} response.Response "获取激励语成功"
// @Failure 400 {object} response.Response "设备号不能为空 或 其他错误"
// @Router /slogan/GetSlogan/{device_num} [get]
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
	c.JSON(http.StatusOK, response.Response{Message: "获取激励语成功"})
}

// ChangeSlogan 更新激励语
// @Summary 更新激励语
// @Description 根据用户ID修改激励语
// @Tags 标语
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param slogan body request.Slogan true "新的激励语"
// @Success 200 {object} response.Response "激励语更新成功"
// @Failure 400 {object} response.Response "无该用户 或 其他错误"
// @Router /slogan/ChangeSlogan/{user_id} [put]
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
