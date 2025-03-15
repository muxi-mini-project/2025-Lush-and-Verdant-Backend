package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/model"
	"2025-Lush-and-Verdant-Backend/tool"
	"fmt"
	"strconv"
)

type GroupService interface {
	AddGroupMember(group *request.ExecuteGroupMember) error
	CreateGroup(group *request.GroupRequest) (*response.GroupInfo, error)
	UpdateGroup(group *request.GroupRequest) (*response.GroupInfo, error)
	DeleteGroup(groupNum uint, executeId uint) error
	DeleteGroupMember(group *request.ExecuteGroupMember) error
	GetGroupInfo(groupNum uint) (*response.GroupInfo, error)
	GetGroupMemberList(groupNum uint) (*response.Users, error)
	GetGroupList(userId uint) (*response.GroupInfos, error)
	GetTenGroup(pn int) ([]response.GroupInfo, error)
}

type GroupServiceImpl struct {
	Dao dao.GroupDAO
}

func NewGroupServiceImpl(Dao dao.GroupDAO) *GroupServiceImpl {
	return &GroupServiceImpl{
		Dao: Dao,
	}
}

func BindGroup(Group *model.Group, group *request.GroupRequest) {
	Group.Name = group.Name
	Group.Description = group.Description
	Group.Password = group.Password
	Group.IsPublic = group.IsPublic
}

func (gsr *GroupServiceImpl) CreateGroup(group *request.GroupRequest) (*response.GroupInfo, error) {
	var Group model.Group
	BindGroup(&Group, group)
	id, err := strconv.Atoi(group.ExecuteId)
	if err != nil {
		return nil, err
	}
	Group.GroupOwnerId = uint(id)
	err = gsr.Dao.CreteGroup(&Group)
	if err != nil {
		return nil, err
	}
	return &response.GroupInfo{
		ID:          tool.UintToString(Group.ID),
		Name:        Group.Name,
		Description: Group.Description,
		IsPublic:    Group.IsPublic,
		GroupOwner:  tool.UintToString(Group.GroupOwnerId),
	}, nil
}

func (gsr *GroupServiceImpl) UpdateGroup(group *request.GroupRequest) (*response.GroupInfo, error) {
	var Group model.Group
	BindGroup(&Group, group)
	groupNum, err := strconv.Atoi(group.GroupNum)
	if err != nil {
		return nil, fmt.Errorf("群号有问题")
	}
	//绑定群id
	Group.ID = uint(groupNum)

	id, err := strconv.Atoi(group.ExecuteId)
	if err != nil {
		return nil, err
	}
	//先检测权限
	ok := gsr.Dao.CheckGroupOwner(&Group, uint(id))
	if ok {
		err := gsr.Dao.UpdateGroup(&Group)
		if err != nil {
			return nil, err
		}
		return &response.GroupInfo{
			ID:          tool.UintToString(Group.ID),
			Name:        Group.Name,
			Description: Group.Description,
			IsPublic:    Group.IsPublic,
			GroupOwner:  tool.UintToString(Group.GroupOwnerId),
		}, nil
	}
	return nil, fmt.Errorf("没有相应的权限")
}

func (gsr *GroupServiceImpl) DeleteGroup(groupNum uint, executeId uint) error {
	var Group model.Group
	//绑定群id
	Group.ID = groupNum
	//检测权限
	ok := gsr.Dao.CheckGroupOwner(&Group, executeId)
	if ok {
		err := gsr.Dao.DeleteGroup(&Group)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("用户没有权限")
}

func (gsr *GroupServiceImpl) GetGroupInfo(groupNum uint) (*response.GroupInfo, error) {
	var groupInfo response.GroupInfo
	group, err := gsr.Dao.GetGroupInfo(groupNum)
	if err != nil {
		return nil, err
	}
	id := tool.UintToString(group.ID)
	groupInfo.ID = id
	groupInfo.Name = group.Name
	groupInfo.Description = group.Description
	groupInfo.IsPublic = group.IsPublic
	groupInfo.GroupOwner = tool.UintToString(group.GroupOwnerId)
	return &groupInfo, nil
}

func (gsr *GroupServiceImpl) GetGroupMemberList(groupNum uint) (*response.Users, error) {
	//获取一手消息
	nums, users, err := gsr.Dao.GetGroupMemberList(groupNum)
	if err != nil {
		return nil, err
	}

	// 预分配 userList 的容量
	userList := make([]response.User, 0, len(users))

	// 遍历 users，映射到 response.User
	for _, v := range users {
		userList = append(userList, response.User{
			UserName:   v.Username,
			Email:      v.Email,
			Slogan:     v.Slogan,
			GoalPublic: v.GoalPublic,
		})
	}

	// 构造返回结果
	userRes := response.Users{
		Nums:  nums,
		Users: userList,
	}

	return &userRes, nil
}

// GetGroupList 获取群聊信息
func (gsr *GroupServiceImpl) GetGroupList(userId uint) (*response.GroupInfos, error) {
	//获取一手消息
	nums, groups, err := gsr.Dao.GetGroupList(userId)
	if err != nil {
		return nil, err
	}

	groupList := make([]response.GroupInfo, 0, len(groups))
	for _, v := range groups {
		groupList = append(groupList, response.GroupInfo{
			ID:          tool.UintToString(v.ID),
			Name:        v.Name,
			Description: v.Description,
			IsPublic:    v.IsPublic,
			GroupOwner:  tool.UintToString(v.GroupOwnerId),
		})
	}
	// 构造返回结果
	groupsRes := response.GroupInfos{
		Nums:   nums,
		Groups: groupList,
	}

	return &groupsRes, nil
}

func (gsr *GroupServiceImpl) AddGroupMember(group *request.ExecuteGroupMember) error {
	useId, err := tool.StringToUint(group.UserId)
	if err != nil {
		return err
	}
	groupNum, err := tool.StringToUint(group.GroupNum)
	if err != nil {
		return err
	}
	return gsr.Dao.AddGroupMember(useId, groupNum)
}

func (gsr *GroupServiceImpl) DeleteGroupMember(group *request.ExecuteGroupMember) error {
	useId, err := tool.StringToUint(group.UserId)
	if err != nil {
		return err
	}
	groupNum, err := tool.StringToUint(group.GroupNum)
	if err != nil {
		return err
	}
	return gsr.Dao.DeleteGroupMember(useId, groupNum)
}

func (gsr *GroupServiceImpl) GetTenGroup(pn int) ([]response.GroupInfo, error) {
	group, err := gsr.Dao.GetTenGroup(pn)
	if err != nil {
		return nil, err
	}

	groupList := make([]response.GroupInfo, 0, len(group))
	for _, v := range group {
		groupList = append(groupList, response.GroupInfo{
			ID:          tool.UintToString(v.ID),
			Name:        v.Name,
			Description: v.Description,
			IsPublic:    v.IsPublic,
			GroupOwner:  tool.UintToString(v.GroupOwnerId),
		})
	}
	return groupList, nil
}
