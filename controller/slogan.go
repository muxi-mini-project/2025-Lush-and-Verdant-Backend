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
// @Success 200 {object} response.Response{data=response.SloganResponse} "获取激励语成功"
// @Failure 400 {object} response.Response "设备号不能为空 或 其他错误"
// @Router /slogan/GetSlogan/{device_num} [get]
func (uc *SloganController) GetSlogan(c *gin.Context) {
	device := c.Param("device_num")
	if device == "" {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "设备号不能为空"})
		return
	}

	slogan, err := uc.ssr.GetSlogan(device)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "激励语获取失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "获取激励语成功", Data: slogan})
}

// ChangeSlogan 更新激励语
// @Summary 更新激励语
// @Description 根据用户ID修改激励语
// @Tags 标语
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param slogan body request.Slogan true "新的激励语"
// @Success 200 {object} response.Response{data=response.SloganResponse} "激励语更新成功"
// @Failure 400 {object} response.Response "无该用户 或 其他错误"
// @Router /slogan/ChangeSlogan/{user_id} [put]
func (uc *SloganController) ChangeSlogan(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无该用户"})
		return
	}

	var newSlogan request.Slogan
	if err := c.ShouldBindJSON(&newSlogan); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "识别标语失败"})
		return
	}

	//显式转换
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "转换失败"})
		return
	}

	err = uc.ssr.ChangeSlogan(uint(idInt), newSlogan)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "修改激励语失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "激励语更新成功", Data: newSlogan})
}

// SearchSlogan 获取激励语
// @Summary 获取激励语
// @Description 用户获取激励语
// @Tags 标语
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.PostGoalResponse} "获取成功"
// @Failure 400 {object} response.Response "解析失败"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /slogan/SearchSlogan [get]
func (uc *SloganController) SearchSlogan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Response{Code: 401, Message: "用户未授权"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "类型转换失败"})
		return
	}
	userIDUint := uint(userIDInt)

	slogan, err := uc.ssr.SearchSlogan(userIDUint)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "获取激励语失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "获取激励语成功", Data: slogan})
}
