package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type Request struct {
	Newslogan string `json:"newslogan"`
}

type Slogans struct {
	Id     int    `json:"id"`
	Slogan string `json:"slogan"`
}

type Users struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Device_Num string `json:"device_num"`
	Email      string `json:"email"`
	Goal_Pubic string `json:"goal_public"`
	Slogan     string `json:"slogan"`
}

type Newslogan struct {
	Newslogan string `json:"newslogan"`
}

// Change_words 通过更新用户的激励语来修改用户口号
// @Summary 更新用户的口号
// @Description 用户可以通过此接口更新他们的口号
// @Tags 用户
// @Accept json
// @Produce json
// @Param newSlogan body Newslogan true "新的口号内容"
// @Security BearerAuth
// @Success 200 {object} Response{message=string} "更新成功"
// @Failure 400 {object} Response{error=string} "请求错误"
// @Failure 404 {object} Response{error=string} "未找到"
// @Failure 500 {object} Response{error=string} "服务器错误"
// @Router /courage_words/change_word/{id} [put]
func Change_words(c *gin.Context) {
	id := c.Param("user_id")

	var newslogan Newslogan
	if err := c.ShouldBind(&newslogan); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	_, err := db.Exec("UPDATE user SET slogan = ? WHERE id = ?", newslogan.Newslogan, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "更新失败"})
		fmt.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, Response{Message: "更新成功"})
}

// Get_words 根据设备号获取激励语
// @Summary 获取设备的激励语
// @Description 通过设备号获取激励语，如果设备尚未设置激励语，则为其分配一个
// @Tags 用户
// @Accept json
// @Produce json
// @Param device path string true "设备号"
// @Success 200 {object} Response{message=string} "更新成功"
// @Failure 400 {object} Response{error=string} "设备号不能为空"
// @Failure 404 {object} Response{error=string} "未找到激励语"
// @Failure 409 {object} Response{message=string} "已拥有激励语"
// @Failure 500 {object} Response{error=string} "服务器错误"
// @Router /courage_words/get_word/{device} [get]
func Get_words(c *gin.Context) {
	device := c.Param("device")

	if device == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "设备号不能为空"})
		return
	}

	var exitslogan string
	err := db.QueryRow("SELECT slogan FROM user WHERE device = ?", device).Scan(&exitslogan)
	if err == nil && exitslogan != "" {
		c.JSON(http.StatusAlreadyReported, Response{Message: "已拥有激励语"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(10) + 1

	var slogans Slogans
	err = db.QueryRow("SELECT id,slogan FROM slogan_list WHERE id = ?", id).Scan(&slogans.Id, &slogans.Slogan)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Error: "读取失败"})
		return
	}

	_, err = db.Exec("UPDATE user SET slogan = ? WHERE device_num = ?", slogans.Slogan, device)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "更新失败"})
		return
	}

	c.JSON(http.StatusOK, Response{Message: "更新成功"})
}
