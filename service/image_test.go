package service

import (
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/dao/mocks"
	"2025-Lush-and-Verdant-Backend/model"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestImageServiceImpl_GetUserImage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	// 创建依赖
	mockDao := daomock.NewMockImageDAO(ctl)
	ImageService := ImageServiceImpl{Dao: mockDao}

	tests := []struct {
		name             string
		mockGetUserImage func()
		expectResult     string
		expectedError    error
	}{
		{
			name: "找到用户并获得url",
			mockGetUserImage: func() {
				mockDao.EXPECT().GetUserImage(&model.User{}).Return("test.com", nil)
			},
			expectResult:  "test.com",
			expectedError: nil,
		}, {
			name: "找到用户但是没有获得url",
			mockGetUserImage: func() {
				mockDao.EXPECT().GetUserImage(&model.User{}).Return("", fmt.Errorf("没有头像"))
			},
			expectResult:  "",
			expectedError: fmt.Errorf("没有头像"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock的行为
			tt.mockGetUserImage()

			var user model.User

			//因为没有创建自己的错误类型
			//而是通过fmt.Errorf()动态创建的，因此即使信息相同，但是地址不同
			url, err := ImageService.GetUserImage(&user)
			//t.Logf("%#v", err)
			//t.Logf("%#v", tt.expectedError)
			if url != tt.expectResult {
				t.Errorf("GetUserImage() url = %v, want %v", url, tt.expectResult)
			}
			//因此不能使用errors.Is来比较
			if err == nil && tt.expectedError != nil || err != nil && tt.expectedError == nil {
				t.Errorf("GetUserImage() err = %v, expect %v", err, tt.expectedError)
			} else if err != nil && tt.expectedError != nil && tt.expectedError.Error() != err.Error() {
				t.Errorf("GetUserImage() err = %v, expect %v", err, tt.expectedError)
			}
			//// 使用 errors.Is 来比较错误
			//if !errors.Is(err, tt.expectedError) {
			//	t.Errorf("GetUserImage() error = %v, want %v", err, tt.expectedError)
			//}
		})
	}
}

func TestImageServiceImpl_GetUserAllImage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	// 创建依赖
	mockDao := daomock.NewMockImageDAO(ctl)
	ImageService := ImageServiceImpl{Dao: mockDao}
	tests := []struct {
		name           string
		mockGetUserAll func()
		expectResult   []string
		expectedError  error
	}{
		{
			name: "找到用户并获取到历史头像",
			mockGetUserAll: func() {
				mockDao.EXPECT().GetUserAllImage(&model.User{}).Return([]string{"test.com"}, nil)
			},
			expectResult:  []string{"test.com"},
			expectedError: nil,
		},
		{
			name: "找到用户但没获得历史头像",
			mockGetUserAll: func() {
				mockDao.EXPECT().GetUserAllImage(&model.User{}).Return(make([]string, 0), fmt.Errorf("没有头像"))
			},
			expectResult:  make([]string, 0),
			expectedError: fmt.Errorf("没有头像"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockGetUserAll()
			var user model.User
			images, err := ImageService.GetUserAllImage(&user)
			//因此不能使用errors.Is来比较
			if err == nil && tt.expectedError != nil || err != nil && tt.expectedError == nil {
				t.Errorf("GetUserAllImage() err = %v, expect %v", err, tt.expectedError)
			} else if err != nil && tt.expectedError != nil && tt.expectedError.Error() != err.Error() {
				t.Errorf("GetUserAllImage() err = %v, expect %v", err, tt.expectedError)
			}

			if images == nil { //先检查是否为nil
				if !reflect.DeepEqual([]string{}, tt.expectResult) {
					t.Errorf("GetUserAllImage() result = %v, want %v", images, tt.expectResult)
				}
			} else {
				if !reflect.DeepEqual(images, tt.expectResult) {
					t.Errorf("GetUserAllImage() result = %v, want %v", images, tt.expectResult)
				}
			}

		})
	}
}
func TestImageServiceImpl_UpdateUserImage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	// 创建依赖
	mockDao := daomock.NewMockImageDAO(ctl)
	imageService := ImageServiceImpl{Dao: mockDao}

	tests := []struct {
		name            string
		mockCreateImage func()
		expectedError   error
	}{
		{
			name: "创建图片成功",
			mockCreateImage: func() {
				mockDao.EXPECT().CreateUserImage(&model.UserImage{}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "创建图片失败",
			mockCreateImage: func() {
				mockDao.EXPECT().CreateUserImage(&model.UserImage{}).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockCreateImage()

			var image model.UserImage
			err := imageService.UpdateUserImage(&image)
			if err == nil && tt.expectedError != nil || err != nil && tt.expectedError == nil {
				t.Errorf("UpdateUserImage() err = %v, expect %v", err, tt.expectedError)
			} else if err != nil && tt.expectedError != nil && tt.expectedError.Error() != err.Error() {
				t.Errorf("UpdateUserImage() err = %v, expect %v", err, tt.expectedError)
			}
		})
	}
}

// 函数替换 --> 可以抽象为接口
//var newPutPolicy = uptoken.NewPutPolicy

func TestImageServiceImpl_GetToken(t *testing.T) {
	tasks := []struct {
		name          string
		qny           *config.QiNiuYunConfig
		expectResult  string
		expectedError error
	}{
		{
			name: "配置文件正常",
			qny: &config.QiNiuYunConfig{
				AccessKey:  "AccessKey",
				SecretKey:  "SecretKey",
				BucketName: "BucketName",
				DomainName: "DomainName",
			},
			expectResult:  "test.com",
			expectedError: nil,
		},
		{
			name: "配置文件中bucket不正常读取",
			qny: &config.QiNiuYunConfig{
				AccessKey:  "AccessKey",
				SecretKey:  "SecretKey",
				BucketName: "",
				DomainName: "DomainName",
			},
			expectResult:  "",
			expectedError: errors.New("failed to set put policy: empty bucket name"),
		},
	}

	for _, tt := range tasks {
		t.Run(tt.name, func(t *testing.T) {

			isr := &ImageServiceImpl{
				qny: tt.qny,
			}
			token, err := isr.GetToken()
			if err == nil && tt.expectedError != nil || err != nil && tt.expectedError == nil {
				t.Errorf("GetToken() err = %v, expect %v", err, tt.expectedError)
			} else if err != nil && tt.expectedError != nil && tt.expectedError.Error() != err.Error() {
				t.Errorf("GetToken() err = %v, expect %v", err, tt.expectedError)
			}
			if token == "" {
				if tt.expectResult != "" {
					t.Errorf("GetToken() result = %v, want %v", token, tt.expectResult)
				}
			}

		})
	}
}

func TestImageServiceImpl_GetGroupImage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	// 创建依赖
	mockDao := daomock.NewMockImageDAO(ctl)
	imageService := &ImageServiceImpl{Dao: mockDao}

	tests := []struct {
		name          string
		mockGetGroup  func()
		expectResult  string
		expectedError error
	}{
		{
			name: "获取小组头像成功",
			mockGetGroup: func() {
				mockDao.EXPECT().GetGroupImage(&model.Group{}).Return("test.com", nil)
			},
			expectResult:  "test.com",
			expectedError: nil,
		},
		{
			name: "获取头像失败",
			mockGetGroup: func() {
				mockDao.EXPECT().GetGroupImage(&model.Group{}).Return("", fmt.Errorf("没有头像"))
			},
			expectResult:  "",
			expectedError: fmt.Errorf("没有头像"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockGetGroup()
			var group model.Group

			url, err := imageService.GetGroupImage(&group)
			if err == nil && tt.expectedError != nil || err != nil && tt.expectedError == nil {
				t.Errorf("UpdateUserImage() err = %v, expect %v", err, tt.expectedError)
			} else if err != nil && tt.expectedError != nil && tt.expectedError.Error() != err.Error() {
				t.Errorf("UpdateUserImage() err = %v, expect %v", err, tt.expectedError)
			}
			if url != tt.expectResult {
				t.Errorf("UpdateUserImage() url = %v, want %v", url, tt.expectResult)
			}
		})
	}

}

func TestImageServiceImpl_GetGroupAllImage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockDao := daomock.NewMockImageDAO(ctl)
	ImageService := ImageServiceImpl{Dao: mockDao}
	tests := []struct {
		name            string
		mockGetGroupAll func()
		expectResult    []string
		expectedError   error
	}{
		{
			name: "获取小组历史头像成功",
			mockGetGroupAll: func() {
				mockDao.EXPECT().GetGroupAllImage(&model.Group{}).Return([]string{"test.com"}, nil)
			},
			expectResult:  []string{"test.com"},
			expectedError: nil,
		},
		{
			name: "获取小组历史头像失败",
			mockGetGroupAll: func() {
				mockDao.EXPECT().GetGroupAllImage(&model.Group{}).Return([]string{}, fmt.Errorf("没有头像"))
			},
			expectResult:  []string{},
			expectedError: fmt.Errorf("没有头像"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockGetGroupAll()

			var group model.Group
			images, err := ImageService.GetGroupAllImage(&group)
			if err == nil && tt.expectedError != nil || err != nil && tt.expectedError == nil {
				t.Errorf("UpdateUserImage() err = %v, expect %v", err, tt.expectedError)
			} else if err != nil && tt.expectedError != nil && tt.expectedError.Error() != err.Error() {
				t.Errorf("UpdateUserImage() err = %v, expect %v", err, tt.expectedError)
			}
			if images == nil { //先检查是否为nil
				if !reflect.DeepEqual([]string{}, tt.expectResult) {
					t.Errorf("GetUserAllImage() result = %v, want %v", images, tt.expectResult)
				}
			} else {
				if !reflect.DeepEqual(images, tt.expectResult) {
					t.Errorf("GetUserAllImage() result = %v, want %v", images, tt.expectResult)
				}
			}
		})
	}
}
func TestImageServiceImpl_UpdateGroupImage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockDao := daomock.NewMockImageDAO(ctl)
	imageService := ImageServiceImpl{Dao: mockDao}
	tests := []struct {
		name            string
		mockUpdateGroup func()
		expectedError   error
	}{
		{
			name: "创建图片成功",
			mockUpdateGroup: func() {
				mockDao.EXPECT().CreateGroupImage(&model.GroupImage{}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "创建图片失败",
			mockUpdateGroup: func() {
				mockDao.EXPECT().CreateGroupImage(&model.GroupImage{}).Return(fmt.Errorf("database error"))
			},
			expectedError: fmt.Errorf("database error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUpdateGroup()
			err := imageService.UpdateGroupImage(&model.GroupImage{})
			if err == nil && tt.expectedError != nil || err != nil && tt.expectedError == nil {
				t.Errorf("UpdateUserImage() err = %v, expect %v", err, tt.expectedError)
			} else if err != nil && tt.expectedError != nil && tt.expectedError.Error() != err.Error() {
				t.Errorf("UpdateUserImage() err = %v, expect %v", err, tt.expectedError)
			}
		})
	}
}
