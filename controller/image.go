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

// 获取上传的token
func (ic *ImageController) GetUpToken(c *gin.Context) {
	uptoken, err := ic.isr.GetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Message: "获取成功", Token: uptoken})
}

func (ic *ImageController) GetUserImage(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "id传输错误"})
		return
	}

	var user model.User
	user.ID = uint(id)
	url, err := ic.isr.GetUserImage(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Data: url})
}

func (ic *ImageController) GetUserAllImage(c *gin.Context) {
	//获取id
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "id传输错误"})
		return
	}
	//得到url
	var user model.User
	user.ID = uint(id)
	urls, err := ic.isr.GetUserAllImage(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Data: urls})
}

func (ic *ImageController) UpdateUserImage(c *gin.Context) {
	var imageRequest request.Image
	if err := c.ShouldBindJSON(&imageRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}
	var image model.UserImage
	image.UserID = uint(imageRequest.Id)
	image.Url = imageRequest.Url
	err := ic.isr.UpdateUserImage(&image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Message: "更新成功"})
}

//group

func (ic *ImageController) GetGroupImage(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "id传输错误"})
		return
	}

	var group model.Group
	group.ID = uint(id)
	url, err := ic.isr.GetGroupImage(&group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Data: url})
}
func (ic *ImageController) GetGroupAllImage(c *gin.Context) {
	//获取id
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "id传输错误"})
		return
	}
	//得到url
	var group model.Group
	group.ID = uint(id)
	urls, err := ic.isr.GetGroupAllImage(&group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{Data: urls})
}
func (ic *ImageController) UpdateGroupImage(c *gin.Context) {
	var imageRequest request.Image
	if err := c.ShouldBindJSON(&imageRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}

	var image model.GroupImage
	image.GroupID = uint(imageRequest.Id)
	image.Url = imageRequest.Url
	//查询group，获得group对象

	// todo 检测权限
	err := ic.isr.UpdateGroupImage(&image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Message: "更新成功"})
}
