package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/service"
	"2025-Lush-and-Verdant-Backend/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GroupController struct {
	gsr service.GroupService
}

func NewGroupController(gsr service.GroupService) *GroupController {
	return &GroupController{gsr: gsr}
}

// CreateGroup 创建新群聊
// @Summary 创建新群聊
// @Description 根据请求参数创建一个新的群聊
// @Tags Group
// @Accept json
// @Produce json
// @Param request body request.GroupRequest true "请求参数"
// @Success 201 {object} response.Response{Data=response.GroupInfo} "成功创建群聊"
// @Failure 400 {object} response.Response "请求参数错误或创建失败"
// @Router /group/create [post]
func (gc *GroupController) CreateGroup(c *gin.Context) {
	//不需要查重，群聊的id唯一就行
	var group *request.GroupRequest
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "解析错误"})
		return
	}

	groupRes, err := gc.gsr.CreateGroup(group)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "创建失败"})
		return
	}
	c.JSON(http.StatusCreated, response.Response{Code: 201, Message: "创建成功", Data: groupRes})
}

// UpdateGroup 更新群聊相关信息
// @Summary 更新群聊相关信息
// @Description 根据请求参数更新群聊的信息
// @Tags Group
// @Accept json
// @Produce json
// @Param request body request.GroupRequest true "请求参数"
// @Success 200 {object} response.Response{Data=response.GroupInfo} "成功更新群聊信息"
// @Failure 400 {object} response.Response "请求参数错误或更新失败"
// @Router /group/update [post]
func (gc *GroupController) UpdateGroup(c *gin.Context) { //todo 返回群聊信息
	var group *request.GroupRequest
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "解析失败"})
		return
	}

	groupRes, err := gc.gsr.UpdateGroup(group)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "更新成功", Data: groupRes})
}

// DeleteGroup 解散群聊
// @Summary 解散群聊
// @Description 根据请求参数删除群聊
// @Tags Group
// @Accept json
// @Produce json
// @Param request body request.GroupRequest true "请求参数"
// @Success 200 {object} response.Response "成功解散群聊"
// @Failure 400 {object} response.Response "请求参数错误或解散失败"
// @Router /group/delete [post]
func (gc *GroupController) DeleteGroup(c *gin.Context) {
	var group *request.GroupRequest
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}

	groupNum, err := strconv.Atoi(group.GroupNum)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	id, err := strconv.Atoi(group.ExecuteId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	err = gc.gsr.DeleteGroup(uint(groupNum), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "解散群聊成功"})
}

// GetGroupInfo 可以获取人数，获取群聊名称，获取群号
// @Summary 获取群聊信息
// @Description 根据群号获取群聊的详细信息，包括人数、群聊名称和群号
// @Tags Group
// @Accept json
// @Produce json
// @Param groupNUM path int true "群号"
// @Success 200 {object} response.Response{data=response.GroupInfo} "成功获取群聊信息"
// @Failure 400 {object} response.Response "请求参数错误或获取失败"
// @Router /group/info/{groupNum} [get]
func (gc *GroupController) GetGroupInfo(c *gin.Context) {
	groupNUm := c.Param("groupNum")
	GroupNum, err := strconv.Atoi(groupNUm)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	group, err := gc.gsr.GetGroupInfo(uint(GroupNum))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: http.StatusOK, Message: "获取成功", Data: group})
}

// GetGroupMemberList 获取群人数和群成员
// @Summary 获取群成员列表
// @Description 根据群号获取群成员列表
// @Tags Group
// @Accept json
// @Produce json
// @Param groupNum path int true "群号"
// @Success 200 {object} response.Response{data=response.Users} "成功获取群成员列表"
// @Failure 400 {object} response.Response "请求参数错误或获取失败"
// @Router /group/members/{groupNum} [get]
func (gc *GroupController) GetGroupMemberList(c *gin.Context) {
	groupNum := c.Param("groupNum")
	//通过群号获取用户群人数
	GroupNum, err := tool.StringToUint(groupNum)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	users, err := gc.gsr.GetGroupMemberList(GroupNum)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "查询群成员成功", Data: users})
}

// GetGroupList 通过用户id获取新的自己的小组
// @Summary 通过用户id获取新的自己的小组
// @Description 通过用户id获取新的自己的小组
// @Tags Group
// @Accept json
// @Produce json
// @Param id path int true "用户id"
// @Success 200 {object} response.Response{data=response.GroupInfos} "成功获取群成员列表"
// @Failure 400 {object} response.Response "请求参数错误或获取失败"
// @Router /group/list/{id} [get]
func (gc *GroupController) GetGroupList(c *gin.Context) {
	idStr := c.Param("id")
	id, err := tool.StringToUint(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	groups, err := gc.gsr.GetGroupList(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "查找成功", Data: groups})
}

// AddGroupMember 加入群聊
// @Summary 加入群聊
// @Description 通过用户id和群号加入群聊
// @Tags Group
// @Accept json
// @Produce json
// @Param request body request.ExecuteGroupMember true "请求参数"
// @Success 200 {object} response.Response "成功加入群聊"
// @Failure 400 {object} response.Response "请求参数错误或加入失败"
// @Router /group/member/add [post]
func (gc *GroupController) AddGroupMember(c *gin.Context) {
	//通过用户id 和群号添加成员
	var addMember *request.ExecuteGroupMember
	if err := c.ShouldBindJSON(&addMember); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	err := gc.gsr.AddGroupMember(addMember)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "用户添加成功"})
}

// DeleteGroupMember 退出群聊
// @Summary 退出群聊
// @Description 通过用户id和群号退出群聊
// @Tags Group
// @Accept json
// @Produce json
// @Param request body request.ExecuteGroupMember true "请求参数"
// @Success 200 {object} response.Response "成功退出群聊"
// @Failure 400 {object} response.Response "请求参数错误或退出失败"
// @Router /group/member/delete [post]
func (gc *GroupController) DeleteGroupMember(c *gin.Context) {
	//通过用户id 和群号退出群聊
	var deleteMember *request.ExecuteGroupMember
	if err := c.ShouldBindJSON(&deleteMember); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	err := gc.gsr.DeleteGroupMember(deleteMember)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "退出群聊成功"})
}

// GetTenGroup 分页获取小组十条信息
// @Summary 分页获取小组十条信息
// @Description 通过页码获取小组前十条信息
// @Tags Group
// @Accept json
// @Produce json
// @Param pn query int false "页码"
// @Success 200 {object} response.Response{data=response.GroupInfos} "成功获取小组前十条信息"
// @Failure 400 {object} response.Response "请求参数错误或获取失败"
// @Router /group/ten [get]
func (gc *GroupController) GetTenGroup(c *gin.Context) {
	pnStr := c.Query("pn")
	if pnStr == "" {
		pnStr = "0"
	}

	pn, err := strconv.Atoi(pnStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	groups, err := gc.gsr.GetTenGroup(pn)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "获取成功", Data: groups})
}
