package dao

import (
	"2025-Lush-and-Verdant-Backend/model"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
)

type GroupDAO interface {
	CreteGroup(group *model.Group) error
	UpdateGroup(group *model.Group) error
	DeleteGroup(group *model.Group) error
	CheckGroupOwner(group *model.Group, useId uint) bool
	GetGroupInfo(groupNum uint) (*model.Group, error)
	GetGroupMemberList(groupNum uint) (int, []model.User, error)
	GetGroupMemberIdList(groupNum uint) ([]int, error)
	GetGroupList(userId uint) (int, []model.Group, error)
	GetGroupIdList(userId uint) ([]int, error)
	AddGroupMember(userId uint, groupNum uint) error
	DeleteGroupMember(userId uint, groupNum uint) error
	GetTenGroup(offset int) ([]model.Group, error)
}

type GroupDAOImpl struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewGroupDAOImpl(db *gorm.DB, rdb *redis.Client) *GroupDAOImpl {
	return &GroupDAOImpl{
		db:  db,
		rdb: rdb,
	}
}

// CreteGroup 创建群聊
// 其中mysql中存储的主要是群的相关介绍
func (dao *GroupDAOImpl) CreteGroup(group *model.Group) error {
	result := dao.db.Create(group)
	if result.Error != nil {
		return fmt.Errorf("创建群聊%s失败", group.Name)
	}

	//绑定群主
	var user model.User
	result = dao.db.Where("id = ?", group.GroupOwnerId).First(&user)
	if result.Error != nil {
		return fmt.Errorf("创建群聊，绑定群主失败")
	}
	err := dao.db.Model(&group).Association("Users").Append(&user)
	if err != nil {
		return fmt.Errorf("创建群聊，绑定群主失败")
	}
	return nil
	// todo 待优化
	// todo 加上缓存
	////在redis中set存取群id和用户id
	//id := strconv.Itoa(int(group.ID))
	//err := dao.rdb.SAdd(context.Background(), "group:members:"+id, group.GroupOwnerId).Err()
	//if err != nil {
	//	return fmt.Errorf("创建群聊，绑定群主失败")
	//}
}

// UpdateGroup 更新群的相关介绍
func (dao *GroupDAOImpl) UpdateGroup(group *model.Group) error {
	result := dao.db.Save(group)
	if result.Error != nil {
		return fmt.Errorf("更新群聊%s信息失败", group.Name)
	}
	return nil
}

// DeleteGroup 解散群聊
// todo 待优化
// todo 同时清除redis的缓存
func (dao *GroupDAOImpl) DeleteGroup(group *model.Group) error {
	result := dao.db.Model(&group).Unscoped().Delete(&group)
	if result.Error != nil {
		return fmt.Errorf("解散群聊%s失败", group.Name)
	}

	//清除群聊的聊天记录
	_, err := dao.rdb.Del(context.TODO(), "group:msg:"+strconv.Itoa(int(group.ID))).Result()
	if err != nil {
		return fmt.Errorf("解散群聊，清除redis聊天记录失败")
	}
	return nil
}

// CheckGroupOwner 检测是否为群主
func (dao *GroupDAOImpl) CheckGroupOwner(group *model.Group, useId uint) bool {

	result := dao.db.Where("id = ?", group.ID).First(&group)
	if result.Error != nil {
		return false
	}
	if group.GroupOwnerId == useId {
		return true
	}
	return false
}

// GetGroupInfo 通过群号获取群聊的基本信息
func (dao *GroupDAOImpl) GetGroupInfo(groupNum uint) (*model.Group, error) {
	var group model.Group
	result := dao.db.Where("id = ?", groupNum).First(&group)
	if result.Error != nil {
		return nil, fmt.Errorf("获取失败")
	}
	return &group, nil
}

// GetGroupMemberList 通过群号获取群聊人数和群成员
func (dao *GroupDAOImpl) GetGroupMemberList(groupNum uint) (int, []model.User, error) {
	var group model.Group
	result := dao.db.Model(&group).Where("id = ?", groupNum).Preload("Users").First(&group)
	if result.Error != nil {
		return 0, nil, fmt.Errorf("查询群成员失败")
	}
	return len(group.Users), group.Users, nil
}

func (dao *GroupDAOImpl) GetGroupMemberIdList(groupNum uint) ([]int, error) {
	nums, groups, err := dao.GetGroupList(groupNum)
	if err != nil {
		return nil, err
	}
	result := make([]int, 0, nums)
	for _, v := range groups {
		result = append(result, int(v.ID))
	}
	return result, nil
}

// GetGroupList 通过用户id查找自己的群数和群信息
func (dao *GroupDAOImpl) GetGroupList(userId uint) (int, []model.Group, error) {
	var user model.User
	result := dao.db.Model(&user).Preload("Groups").Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return 0, nil, fmt.Errorf("查询群聊失败")
	}
	return len(user.Groups), user.Groups, nil
}

// GetGroupIdList 通过用户id获取群聊id列表
func (dao *GroupDAOImpl) GetGroupIdList(userId uint) ([]int, error) {
	nums, groups, err := dao.GetGroupList(userId)
	if err != nil {
		return nil, err
	}
	result := make([]int, 0, nums)
	for _, v := range groups {
		result = append(result, int(v.ID))
	}
	return result, nil
}

func (dao *GroupDAOImpl) AddGroupMember(userId uint, groupNum uint) error {
	var user model.User
	var group model.Group
	err := dao.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return fmt.Errorf("查询用户失败")
	}
	err = dao.db.Where("id = ?", groupNum).First(&group).Error
	if err != nil {
		return fmt.Errorf("查询群聊失败")
	}
	err = dao.db.Model(&group).Association("Users").Append(&user)
	if err != nil {
		return fmt.Errorf("添加用户失败")
	}
	return nil
}

func (dao *GroupDAOImpl) DeleteGroupMember(userId uint, groupNum uint) error {
	var user model.User
	var group model.Group
	err := dao.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return fmt.Errorf("查询用户失败")
	}
	err = dao.db.Where("id = ?", groupNum).First(&group).Error
	if err != nil {
		return fmt.Errorf("查询群聊失败")
	}
	err = dao.db.Model(&group).Association("Users").Delete(&user)
	if err != nil {
		return fmt.Errorf("退出群聊失败")
	}
	return nil
}

// GetTenGroup 分页获取消息
func (dao *GroupDAOImpl) GetTenGroup(offset int) ([]model.Group, error) {
	var groups []model.Group
	result := dao.db.Limit(10).Offset(offset).Find(&groups)
	if result.Error != nil {
		return nil, fmt.Errorf("获取小组失败")
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("没有更多小组了")
	}
	return groups, nil
}
