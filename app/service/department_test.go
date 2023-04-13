package service

import (
	"asset-management/app/define"
	"testing"

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

}
