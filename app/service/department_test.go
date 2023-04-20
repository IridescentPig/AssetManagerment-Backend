package service

import (
	"asset-management/app/define"
	"asset-management/utils"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDepartment(t *testing.T) {
	InitForTest()

	_ = EntityService.CreateEntity("test_entity")
	//assert.Equal(t, nil, err, "service error")

	err := DepartmentService.CreateDepartment("test_department", 1, 0)
	assert.Equal(t, nil, err, "service error")

	department, err := DepartmentService.GetDepartmentInfoByID(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, "test_department", department.Name, "service error")

	err = DepartmentService.CreateDepartmentUser(define.CreateDepartmentUserReq{
		UserName:        "test_user",
		Password:        "123456",
		DepartmentSuper: false,
	}, 1, 1)
	assert.Equal(t, nil, err, "service error")

	users, err := DepartmentService.GetAllUsers(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, 1, len(users), "service error")

	err = DepartmentService.CreateDepartment("sub_department", 1, 1)
	assert.Equal(t, nil, err, "service error")
	subs, err := DepartmentService.GetSubDepartments(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, 1, len(subs), "service error")

	err = DepartmentService.SetDepartmentManager("test_user", 1)
	assert.Equal(t, nil, err, "service error")

	list, err := DepartmentService.GetDepartmentManagerList(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, 1, len(list), "service error")

	exist, err := DepartmentService.DepartmentHasUsers(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, true, exist, "service error")

	err = DepartmentService.DeleteDepartmentManager(1)
	assert.Equal(t, nil, err, "service error")

	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)
	context := utils.Context{Context: ctx}
	ctx.Set("user", define.UserBasicInfo{
		UserID:          1,
		UserName:        "admin",
		EntitySuper:     true,
		DepartmentSuper: true,
		SystemSuper:     true,
	})
	DepartmentService.GetHeaderDepartmentID(&context)
	DepartmentService.CheckIsInDepartment(&context, 1)
	DepartmentService.CheckIsAncestor(1, 2)
	DepartmentService.ExistsDepartmentByID(1)
	DepartmentService.ExistsDepartmentSub("sub_department", 1, 1)
	DepartmentService.CheckDepartmentIdentity(&context, 1, 1)
	DepartmentService.GetSubDepartmentTreeNodes(1, 1)
}
