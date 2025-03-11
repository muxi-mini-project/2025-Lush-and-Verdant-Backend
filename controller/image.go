package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/model"
	"2025-Lush-and-Verdant-Backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ImageController struct {
	isr service.ImageService
}

func NewImageController(isr service.ImageService) *ImageController {
	return &ImageController{isr: isr}
}

// GetUpToken 获取上传的token
// @Summary 获取上传的token
// @Description 获取用于上传图片的token
// @Tags image
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.UpToken}
// @Failure 500 {object} response.Response
// @Router /get_token [get]
func (ic *ImageController) GetUpToken(c *gin.Context) {
	uptoken, err := ic.isr.GetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "用户未登录"})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "获取成功", Data: uptoken})
}

// GetUserImage 获取用户图片
// @Summary 获取用户图片
// @Description 根据用户ID获取用户的头像
// @Tags image
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=response.URL}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /image/user/get/{id} [get]
func (ic *ImageController) GetUserImage(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "id传输错误"})
		return
	}

	var user model.User
	user.ID = uint(id)
	url, err := ic.isr.GetUserImage(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "获取头像失败"})
		return
	}
	c.JSON(http.StatusOK, response.Response{Data: url})
}

// GetUserAllImage 获取所有用户图片
// @Summary 获取所有用户图片
// @Description 获取用户上传的所有图片
// @Tags image
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=response.URLs}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /image/user/history/{id} [get]
func (ic *ImageController) GetUserAllImage(c *gin.Context) {
	// 获取id
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "id传输错误"})
		return
	}
	// 得到url
	var user model.User
	user.ID = uint(id)
	urls, err := ic.isr.GetUserAllImage(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "获取历史头像失败"})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Data: urls})
}

// UpdateUserImage 更新用户头像
// @Summary 更新用户头像
// @Description 更新用户的头像
// @Tags image
// @Accept json
// @Produce json
// @Param data body request.Image true "图片请求体"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /image/user/update [put]
func (ic *ImageController) UpdateUserImage(c *gin.Context) {
	var imageRequest request.Image
	if err := c.ShouldBindJSON(&imageRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "解析失败"})
		return
	}
	var image model.UserImage
	id, err := strconv.Atoi(imageRequest.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "解析失败"})
		return
	}
	image.UserID = uint(id)
	image.Url = imageRequest.Url
	err = ic.isr.UpdateUserImage(&image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "头像上传失败"})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "更新成功"})
}

// GetGroupImage 获取群组图片
// @Summary 获取群组图片
// @Description 获取群组的头像
// @Tags image
// @Accept json
// @Produce json
// @Param id path int true "群组ID"
// @Success 200 {object} response.Response{data=response.URL}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /image/group/{id} [get]
func (ic *ImageController) GetGroupImage(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "id传输错误"})
		return
	}

	var group model.Group
	group.ID = uint(id)
	url, err := ic.isr.GetGroupImage(&group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "获取群组头像失败"})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Data: url})
}

// GetGroupAllImage 获取所有群组图片
// @Summary 获取所有群组图片
// @Description 获取群组上传的所有图片
// @Tags image
// @Accept json
// @Produce json
// @Param id path int true "群组ID"
// @Success 200 {object} response.Response{data=response.URLs}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /image/group/history/{id} [get]
func (ic *ImageController) GetGroupAllImage(c *gin.Context) {
	// 获取id
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "id传输错误"})
		return
	}
	// 得到url
	var group model.Group
	group.ID = uint(id)
	urls, err := ic.isr.GetGroupAllImage(&group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "获取群组历史头像失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Data: urls})
}

// UpdateGroupImage 更新群组头像
// @Summary 更新群组头像
// @Description 更新群组的头像
// @Tags image
// @Accept json
// @Produce json
// @Param data body request.Image true "图片请求体"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /image/group/update [put]
func (ic *ImageController) UpdateGroupImage(c *gin.Context) {
	var imageRequest request.Image
	if err := c.ShouldBindJSON(&imageRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "解析失败"})
		return
	}

	var image model.GroupImage
	id, err := strconv.Atoi(imageRequest.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "解析失败"})
		return
	}
	image.GroupID = uint(id)
	image.Url = imageRequest.Url
	// 查询group，获得group对象

	// todo 检测权限
	err = ic.isr.UpdateGroupImage(&image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "更新群组头像失败"})
		return
	}
	c.JSON(http.StatusOK, response.Response{Message: "更新成功"})
}
