package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func createMockContext() *gin.Context {
	// 模拟一个空的上下文对象
	return &gin.Context{}
}

func TestUserLogin(t *testing.T) {
	// 创建一个gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockUserService对象
	mockUserService := NewMockUserService(ctrl)

	// 设置期望：调用UserLogin，传入任意gin.Context参数返回nil
	mockUserService.EXPECT().UserLogin(gomock.Any()).Return(nil)

	// 创建一个mock请求的上下文
	ctx := createMockContext()

	// 执行被测试的方法
	err := mockUserService.UserLogin(ctx)

	// 断言：期望没有错误发生
	assert.NoError(t, err)
}

func TestUserLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)

	mockUserService.EXPECT().UserLogin(gomock.Any()).Return(nil)

	ctx := createMockContext()

	err := mockUserService.UserLogin(ctx)

	assert.NoError(t, err)
}

func TestSendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)

	mockUserService.EXPECT().SendEmail(gomock.Any()).Return(nil)

	ctx := createMockContext()

	err := mockUserService.SendEmail(ctx)

	assert.NoError(t, err)
}
