package api

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDepartment(t *testing.T) {
	res := httptest.NewRecorder()
	_, r := gin.CreateTestContext(res)
	InitForTest(r)

	admin := model.User{
		UserName:        "admin",
		Password:        utils.CreateMD5("21232f297a57a5a743894a0e4a801fc3"),
		SystemSuper:     true,
		EntitySuper:     true,
		DepartmentSuper: true,
		Ban:             false,
	}
	dao.UserDao.Create(admin)

	UserLogin := define.UserLoginReq{
		UserName: "admin",
		Password: "21232f297a57a5a743894a0e4a801fc3",
	}
	{
		req := GetRequest(http.MethodPost, "/user/login", headerJson, GetJsonBody(UserLogin))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	b, err := io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]interface{}{}
	json.Unmarshal(b, &data)
	user := data["data"].(map[string]interface{})
	token := user["token"].(string)
	headerJsonToken["Authorization"] = token
	headerFormToken["Authorization"] = token

	CreateEntity := define.CreateEntityReq{
		EntityName: "test_entity1",
	}
	{
		req := GetRequest(http.MethodPost, "/entity/", headerFormToken, GetJsonBody(CreateEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/", headerJson, GetJsonBody(CreateEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	password := "123456"
	managerReq2 := define.ManagerReq{
		Username: "entity_manager",
		Password: &password,
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerFormToken, GetJsonBody(managerReq2))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerJson, GetJsonBody(managerReq2))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}
	UserLogin2 := define.UserLoginReq{
		UserName: "entity_manager",
		Password: "123456",
	}
	{
		req := GetRequest(http.MethodPost, "/user/login", headerJson, GetJsonBody(UserLogin2))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	b, err = io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}
	data = map[string]interface{}{}
	json.Unmarshal(b, &data)
	user = data["data"].(map[string]interface{})
	token = user["token"].(string)
	headerJsonToken["Authorization"] = token
	headerFormToken["Authorization"] = token

	// POST /:entity_id/department
	CreateDepartment := define.CreateDepartmentReq{
		DepartmentName: "test_department",
	}

	{
		req := GetRequest(http.MethodPost, "/entity/1/department", headerFormToken, GetJsonBody(CreateDepartment))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department", headerJson, GetJsonBody(CreateDepartment))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/list
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/list", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/:department_id
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	CreateDepartmentUser := define.CreateDepartmentUserReq{
		UserName:        "test_user",
		Password:        "123456",
		DepartmentSuper: true,
	}

	// POST /:entity_id/department/:department_id/user
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/user", headerFormToken, GetJsonBody(CreateDepartmentUser))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/user", headerJson, GetJsonBody(CreateDepartmentUser))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/user", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	// POST /:entity_id/department/:department_id/manager
	SetDepartmentManager := define.SetDepartmentManagerReq{
		UserName: "test_user",
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/manager", headerFormToken, GetJsonBody(SetDepartmentManager))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/manager", headerJson, GetJsonBody(SetDepartmentManager))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/manager", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/:department_id/user/list
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1/user/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1/user/list", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/:department_id/manager
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1/manager", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1/manager", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// POST /:entity_id/department/:department_id/department
	CreateDepartment2 := define.CreateDepartmentReq{
		DepartmentName: "test_sub_department",
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/department", headerFormToken, GetJsonBody(CreateDepartment2))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/department", headerJson, GetJsonBody(CreateDepartment2))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/department", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/:department_id/department/list
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1/department/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1/department/list", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/tree
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/tree", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/tree", headerForm, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// DELETE /:entity_id/department/:department_id/manager/:user_id
	{
		req := GetRequest(http.MethodDelete, "/entity/1/department/1/manager/3", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodDelete, "/entity/1/department/1/manager/3", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// DELETE /:entity_id/department/:department_id
	{
		req := GetRequest(http.MethodDelete, "/entity/1/department/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodDelete, "/entity/1/department/1", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

}

func TestDepartmentNoPermission(t *testing.T) {
	res := httptest.NewRecorder()
	_, r := gin.CreateTestContext(res)
	InitForTest(r)

	admin := model.User{
		UserName:        "admin",
		Password:        utils.CreateMD5("21232f297a57a5a743894a0e4a801fc3"),
		SystemSuper:     true,
		EntitySuper:     true,
		DepartmentSuper: true,
		Ban:             false,
	}
	dao.UserDao.Create(admin)

	UserLogin := define.UserLoginReq{
		UserName: "admin",
		Password: "21232f297a57a5a743894a0e4a801fc3",
	}
	{
		req := GetRequest(http.MethodPost, "/user/login", headerJson, GetJsonBody(UserLogin))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	b, err := io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]interface{}{}
	json.Unmarshal(b, &data)
	user := data["data"].(map[string]interface{})
	token := user["token"].(string)
	headerJsonToken["Authorization"] = token
	headerFormToken["Authorization"] = token

	CreateEntity := define.CreateEntityReq{
		EntityName: "test_entity1",
	}
	{
		req := GetRequest(http.MethodPost, "/entity/", headerFormToken, GetJsonBody(CreateEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	password := "123456"
	managerReq2 := define.ManagerReq{
		Username: "entity_manager",
		Password: &password,
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerFormToken, GetJsonBody(managerReq2))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	UserLogin2 := define.UserLoginReq{
		UserName: "entity_manager",
		Password: "123456",
	}
	{
		req := GetRequest(http.MethodPost, "/user/login", headerJson, GetJsonBody(UserLogin2))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	b, err = io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}
	data = map[string]interface{}{}
	json.Unmarshal(b, &data)
	user = data["data"].(map[string]interface{})
	token = user["token"].(string)
	headerJsonToken["Authorization"] = token
	headerFormToken["Authorization"] = token

	// POST /:entity_id/department
	CreateDepartment := define.CreateDepartmentReq{
		DepartmentName: "test_department",
	}

	{
		req := GetRequest(http.MethodPost, "/entity/1/department", headerFormToken, GetJsonBody(CreateDepartment))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	no := model.User{
		UserName: "no",
		Password: utils.CreateMD5("21232f297a57a5a743894a0e4a801fc3"),
		Ban:      false,
	}
	dao.UserDao.Create(no)

	UserLogin = define.UserLoginReq{
		UserName: "no",
		Password: "21232f297a57a5a743894a0e4a801fc3",
	}
	{
		req := GetRequest(http.MethodPost, "/user/login", headerJson, GetJsonBody(UserLogin))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	b, err = io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}
	data = map[string]interface{}{}
	json.Unmarshal(b, &data)
	user = data["data"].(map[string]interface{})
	token = user["token"].(string)
	headerJsonToken["Authorization"] = token
	headerFormToken["Authorization"] = token

	CreateEntity = define.CreateEntityReq{
		EntityName: "test_entity1",
	}
	{
		req := GetRequest(http.MethodPost, "/entity/", headerFormToken, GetJsonBody(CreateEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	password = "123456"
	managerReq2 = define.ManagerReq{
		Username: "entity_manager",
		Password: &password,
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerFormToken, GetJsonBody(managerReq2))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// POST /:entity_id/department
	CreateDepartment = define.CreateDepartmentReq{
		DepartmentName: "test_department",
	}

	{
		req := GetRequest(http.MethodPost, "/entity/1/department", headerFormToken, GetJsonBody(CreateDepartment))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/list
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/:department_id
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	CreateDepartmentUser := define.CreateDepartmentUserReq{
		UserName:        "test_user",
		Password:        "123456",
		DepartmentSuper: true,
	}

	// POST /:entity_id/department/:department_id/user
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/user", headerFormToken, GetJsonBody(CreateDepartmentUser))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// POST /:entity_id/department/:department_id/manager
	SetDepartmentManager := define.SetDepartmentManagerReq{
		UserName: "test_user",
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/manager", headerFormToken, GetJsonBody(SetDepartmentManager))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/:department_id/user/list
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1/user/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/:department_id/manager
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1/manager", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/tree
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/tree", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// POST /:entity_id/department/:department_id/department
	CreateDepartment2 := define.CreateDepartmentReq{
		DepartmentName: "test_sub_department",
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/department", headerFormToken, GetJsonBody(CreateDepartment2))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// GET /:entity_id/department/:department_id/department/list
	{
		req := GetRequest(http.MethodGet, "/entity/1/department/1/department/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// DELETE /:entity_id/department/:department_id/manager/:user_id
	{
		req := GetRequest(http.MethodDelete, "/entity/1/department/1/manager/3", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// DELETE /:entity_id/department/:department_id
	{
		req := GetRequest(http.MethodDelete, "/entity/1/department/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}
}
