package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LikeController struct {
	lsr service.LikeService
}

func NewLikeController(lsr service.LikeService) *LikeController {
	return &LikeController{
		lsr: lsr,
	}
}

// Like 点赞和取消点赞
// @Summary 点赞和取消点赞
// @Description 点赞和取消点赞
// @Tags Like
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "点赞请求成功"
// @Failure 400 {object} response.Response "请求参数错误或获取失败"
// @Router /like/send [post]
func (lc *LikeController) Like(c *gin.Context) {
	var like request.ForestLikeReq
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "解析信息失败",
		})
		return
	}
	err := lc.lsr.SendMsg(&like)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code:    500,
			Message: err.Error(),
		})
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "点赞请求成功",
	})
}

// GetForestAllLikes 获取森林点赞数
// @Summary 获取森林点赞数
// @Description 获取森林点赞数
// @Tags Like
// @Param id path int true "用户id"
// @Produce json
// @Success 200 {object} response.Response{Data=response.Likes} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误或获取失败"
// @Router /like/get/{to} [get]
func (lc *LikeController) GetForestAllLikes(c *gin.Context) {
	to := c.Param("to")

	nums, err := lc.lsr.GetLikes(to)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}
	likes := response.Likes{
		NUms: nums,
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "获取成功",
		Data:    likes,
	})
}

// GetForestLikeStatus 获取点赞状态
// @Summary 获取点赞状态
// @Description 获取点赞状态
// @Tags Like
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "已点赞"
// @Failure 400 {object} response.Response "请求参数错误或获取失败"
// @Router /like/status [post]
func (lc *LikeController) GetForestLikeStatus(c *gin.Context) {
	var like request.ForestLikeReq
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	ok := lc.lsr.GetForestLikeStatus(&like)
	if !ok {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "未点赞",
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "已点赞",
	})
}
