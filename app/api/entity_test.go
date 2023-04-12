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

func TestEntity(t *testing.T) {
	res := httptest.NewRecorder()
	_, r := gin.CreateTestContext(res)
	Init(r)

	admin := model.User{
		UserName:    "admin",
		Password:    utils.CreateMD5("21232f297a57a5a743894a0e4a801fc3"),
		SystemSuper: true,
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
		EntityName: "test_entity1",
	}

	{
		req := GetRequest(http.MethodPost, "/entity/", headerJsonToken, GetJsonBody(CreateEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		//b, err := io.ReadAll(res.Result().Body)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//data := map[string]interface{}{}
		//json.Unmarshal(b, &data)
		//log.Printf(((data["error"].(map[string]interface{}))["error"].(map[string]interface{}))[""])

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodGet, "/entity/list", headerJsonToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
		/*b, err := io.ReadAll(res.Result().Body)
		if err != nil {
			log.Fatal(err)
		}
		data := map[string]interface{}{}
		json.Unmarshal(b, &data)
		data = data["data"].(map[string]interface{})
		entity_list := data["entity_list"].([]model.Entity)
		assert.Equal(t, 1, len(entity_list), "response failed")
		assert.Equal(t, "test_entity", entity_list[0], "response failed")*/
	}

	name := "new_entity_name"
	des := "description"

	modifyEntityInfoReq := define.ModifyEntityInfoReq{
		EntityName:  &name,
		Description: &des,
	}

	{
		req := GetRequest(http.MethodPatch, "/entity/1", headerJsonToken, GetJsonBody(modifyEntityInfoReq))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodGet, "/entity/1", headerJsonToken, nil)
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

	managerReq1 := define.ManagerReq{
		Username: "admin",
	}

	password := "123456"
	managerReq2 := define.ManagerReq{
		Username: "entity_manager",
		Password: &password,
	}

	{
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerJsonToken, GetJsonBody(managerReq1))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodPost, "/entity/1/manager", headerJsonToken, GetJsonBody(managerReq2))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

}
