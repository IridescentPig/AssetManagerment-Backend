package service

import (
	"asset-management/app/define"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	InitForTest()

	err := EntityService.CreateEntity("test_entity222")
	assert.Equal(t, nil, err, "service error")

	entities, err := EntityService.GetAllEntity()
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, 2, len(entities), "service error")
	assert.Equal(t, "test_entity", entities[0].Name, "service error")

	exist, err := EntityService.ExistsEntityByName("test_entity222")
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, true, exist, "service error")
	exist, err = EntityService.ExistsEntityByName("not_exist")
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, false, exist, "service error")

	exist, err = EntityService.ExistsEntityByID(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, true, exist, "service error")
	exist, err = EntityService.ExistsEntityByID(3)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, false, exist, "service error")

	entity, err := EntityService.GetEntityInfoByID(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, "test_entity", entity.Name, "service error")

	err = EntityService.CreateManager("manager", "123456", 1)
	assert.Equal(t, nil, err, "service error")

	err = EntityService.SetManager("test_manager")
	assert.Equal(t, nil, err, "service error")

	users, err := EntityService.GetUsersUnderEntity(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, 1, len(users), "service error")

	managers, err := EntityService.GetEntityManagerList(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, 1, len(managers), "service error")

	err = EntityService.DeleteManager(1)
	assert.Equal(t, nil, err, "service error")

	new_name := "new_name"
	des := "description"
	err = EntityService.ModifyEntity(1, define.ModifyEntityInfoReq{
		EntityName:  &new_name,
		Description: &des,
	})
	assert.Equal(t, nil, err, "service error")
	entity, err = EntityService.GetEntityInfoByID(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, new_name, entity.Name, "service error")
	assert.Equal(t, des, entity.Description, "service error")

	exist, err = EntityService.EntityHasUser(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, true, exist, "service error")

	departments, err := EntityService.GetAllDepartmentsUnderEntity(1)
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, 4, len(departments), "service error")

	err = EntityService.DeleteEntity(1)
	assert.Equal(t, nil, err, "service error")

}
