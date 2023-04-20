package api

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/middleware"
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

// assistant function
func print_errormessage(res *httptest.ResponseRecorder) string {
	b, err := io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]interface{}{}
	json.Unmarshal(b, &data)
	data = data["error"].(map[string]interface{})
	code := data["code"].(float64)
	msg := data["message"].(string)
	log.Print("code ", code, ";message ", msg)
	return msg
}

func InitForEntity(r *gin.Engine) {
	group := r.Group("/entity")
	group.Use(utils.Handler(middleware.JWTMiddleware()))
	group.GET("/:entity_id/user/list", utils.Handler(EntityApi.UsersInEntity))             //
	group.GET("/:entity_id/department/list", utils.Handler(EntityApi.DepartmentsInEntity)) // change later
	group.PATCH("/:entity_id", utils.Handler(EntityApi.ModifyEntityInfo))                  //

	group.POST("/:entity_id/department", utils.Handler(DepartmentApi.CreateDepartment))                                          //
	group.POST("/:entity_id/department/:department_id/department", utils.Handler(DepartmentApi.CreateDepartment))                //
	group.DELETE("/:entity_id/department/:department_id", utils.Handler(DepartmentApi.DeleteDepartment))                         //
	group.GET("/:entity_id/department/:department_id", utils.Handler(DepartmentApi.GetDepartmentByID))                           //
	group.GET("/:entity_id/department/:department_id/department/list", utils.Handler(DepartmentApi.GetSubDepartments))           //
	group.GET("/:entity_id/department/:department_id/user/list", utils.Handler(DepartmentApi.GetAllUsersUnderDepartment))        //
	group.POST("/:entity_id/department/:department_id/user", utils.Handler(DepartmentApi.CreateUserInDepartment))                //
	group.POST("/:entity_id/department/:department_id/manager", utils.Handler(DepartmentApi.SetManager))                         //
	group.DELETE("/:entity_id/department/:department_id/manager/:user_id", utils.Handler(DepartmentApi.DeleteDepartmentManager)) //
	group.GET("/:entity_id/department/:department_id/manager", utils.Handler(DepartmentApi.GetDepartmentManager))                //
	group.GET("/:entity_id/department/tree", utils.Handler(DepartmentApi.GetDepartmentTree))

	group.Use(utils.Handler(middleware.CheckSystemSuper()))
	{
		group.POST("/", utils.Handler(EntityApi.CreateEntity))                               //
		group.DELETE("/:entity_id", utils.Handler(EntityApi.DeleteEntity))                   //
		group.GET("/list", utils.Handler(EntityApi.GetEntityList))                           //
		group.GET("/:entity_id", utils.Handler(EntityApi.GetEntityByID))                     //
		group.POST("/:entity_id/manager", utils.Handler(EntityApi.SetManager))               //
		group.DELETE("/:entity_id/manager/:user_id", utils.Handler(EntityApi.DeleteManager)) //
	}
}

func TestEntity(t *testing.T) {
	res := httptest.NewRecorder()
	_, r := gin.CreateTestContext(res)
	InitForTest(r)

	admin := model.User{
		UserName:    "admin",
		Password:    utils.CreateMD5("21232f297a57a5a743894a0e4a801fc3"),
		SystemSuper: true,
		EntitySuper: true,
		Ban:         false,
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
		EntityName: "test_entity111",
	}

	// POST /entity/
	{
		req := GetRequest(http.MethodPost, "/entity/", headerFormToken, GetJsonBody(CreateEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		print_errormessage(res)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/", headerJson, GetJsonBody(CreateEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		//print_errormessage(res)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		//print_errormessage(res)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	// GET /entity/list
	{
		req := GetRequest(http.MethodGet, "/entity/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/entity/list", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// PATCH /entity/1
	name := "new_entity_name"
	des := "description"

	modifyEntityInfoReq := define.ModifyEntityInfoReq{
		EntityName:  &name,
		Description: &des,
	}

	{
		req := GetRequest(http.MethodPatch, "/entity/1", headerFormToken, GetJsonBody(modifyEntityInfoReq))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPatch, "/entity/1", headerJson, GetJsonBody(modifyEntityInfoReq))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPatch, "/entity/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	// GET /entity/1
	{
		req := GetRequest(http.MethodGet, "/entity/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
		b, err := io.ReadAll(res.Result().Body)
		if err != nil {
			log.Fatal(err)
		}
		data := map[string]interface{}{}
		json.Unmarshal(b, &data)
		data = data["data"].(map[string]interface{})
		assert.Equal(t, name, data["entity_name"].(string), "response failed")
		assert.Equal(t, des, data["description"].(string), "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/entity/1", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// POST /entity/1/manager
	managerReq1 := define.ManagerReq{
		Username: "admin",
	}

	{
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerFormToken, GetJsonBody(managerReq1))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerJson, GetJsonBody(managerReq1))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
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
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerForm, GetJsonBody(managerReq2))
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

	// GET /entity/{entity_id}/user/list
	{
		req := GetRequest(http.MethodGet, "/entity/1/user/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/entity/1/user/list", headerForm, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// DELETE /entity/:entity_id/manager/:user_id
	{
		req := GetRequest(http.MethodDelete, "/entity/1/manager/2", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodDelete, "/entity/1/manager/2", headerForm, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// DELETE /entity/:entity_id
	{
		req := GetRequest(http.MethodDelete, "/entity/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodDelete, "/entity/1", headerForm, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

}

func TestEntityNoPermission(t *testing.T) {
	res := httptest.NewRecorder()
	_, r := gin.CreateTestContext(res)
	InitForTest(r)

	admin := model.User{
		UserName: "no",
		Password: utils.CreateMD5("21232f297a57a5a743894a0e4a801fc3"),
		Ban:      false,
	}
	dao.UserDao.Create(admin)

	UserLogin := define.UserLoginReq{
		UserName: "no",
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
		EntityName: "test_entity111",
	}

	// POST /entity/
	{
		req := GetRequest(http.MethodPost, "/entity/", headerFormToken, GetJsonBody(CreateEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		//print_errormessage(res)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// GET /entity/list
	{
		req := GetRequest(http.MethodGet, "/entity/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// PATCH /entity/1
	name := "new_entity_name"
	des := "description"

	modifyEntityInfoReq := define.ModifyEntityInfoReq{
		EntityName:  &name,
		Description: &des,
	}

	{
		req := GetRequest(http.MethodPatch, "/entity/1", headerFormToken, GetJsonBody(modifyEntityInfoReq))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// GET /entity/1
	{
		req := GetRequest(http.MethodGet, "/entity/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// POST /entity/1/manager
	managerReq1 := define.ManagerReq{
		Username: "admin",
	}

	{
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerFormToken, GetJsonBody(managerReq1))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
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
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// GET /entity/{entity_id}/user/list
	{
		req := GetRequest(http.MethodGet, "/entity/1/user/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// DELETE /entity/:entity_id/manager/:user_id
	{
		req := GetRequest(http.MethodDelete, "/entity/1/manager/2", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}

	// DELETE /entity/:entity_id
	{
		req := GetRequest(http.MethodDelete, "/entity/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
	}
}
