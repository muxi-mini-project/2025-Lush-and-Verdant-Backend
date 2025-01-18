package controller

import (
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SloganController struct {
}

func NewSloganController() *SloganController {
	return &SloganController{}
}

func (uc *SloganController) GetSlogan(c *gin.Context) {
	device := c.Param("device_num")

	if device == "" {
		c.JSON(http.StatusBadRequest, model.Response{Error: "设备号不能为空"})
		return
	}

	db := dao.NewDB(dsn)

	db.AutoMigrate(&model.Slogan{})

	// 查找所有的激励语
	var slogans []model.Slogan
	if err := db.Find(&slogans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Error: "无法获取激励语"})
		return
	}

	// 如果没有可用的激励语，返回错误
	if len(slogans) == 0 {
		c.JSON(http.StatusBadRequest, model.Response{Error: "无可用激励语"})
		return
	}

	var slogan model.Slogan
	db.Table("slogans").Order("RAND()").First(&slogan) // 使用Order随机排序，选第一条slogan

	var user model.User
	err := db.Table("users").Where("device_num = ?", device).First(&user)
	if err.Error != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: "找不到对应的设备号用户"})
		return
	}

	user.Slogan = slogan.Slogan
	err = db.Table("users").Save(&user)
	if err.Error != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: "更新激励语失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Message: "获取激励语成功"})
}

func (uc *SloganController) ChangeSlogan(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.Response{Error: "无该用户"})
		return
	}
	db := dao.NewDB(dsn)
	db.AutoMigrate(&model.Slogan{})
	db.AutoMigrate(&model.User{})
	var newslogan model.Slogan
	if err := c.ShouldBind(&newslogan); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: "识别标语失败"})
		return
	}
	var user model.User
	err := db.Table("users").Where("id = ?", id).First(&user).Error
	if err.Error != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: "未找到相关用户"})
		return
	}
	user.Slogan = newslogan.Slogan
	err = db.Table("users").Save(&user).Error
	if err.Error != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response{Message: "激励语更新成功"})
}
