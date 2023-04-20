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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func InitForAsset(r *gin.Engine) {
	group := r.Group("/department")
	group.Use(utils.Handler(middleware.JWTMiddleware()))
	group.GET("/:department_id/asset/list", utils.Handler(AssetApi.GetAssetList))
	group.PATCH("/:department_id/asset/:asset_id", utils.Handler(AssetApi.ModifyAssetInfo))
	group.POST("/:department_id/asset", utils.Handler(AssetApi.CreateAssets))
	group.PATCH("/:department_id/asset/expire", utils.Handler(AssetApi.ExpireAsset))
	group.POST("/:department_id/asset/transfer", utils.Handler(AssetApi.TransferAssets))
}

func TestAsset(t *testing.T) {
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
		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
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
		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	SetDepartmentManager := define.SetDepartmentManagerReq{
		UserName: "entity_manager",
	}
	{
		req := GetRequest(http.MethodPost, "/entity/1/department/1/manager", headerFormToken, GetJsonBody(SetDepartmentManager))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
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

	// POST /:department_id/asset_class
	CreateAssetClass := define.CreateAssetClassReq{
		ClassName: "okok",
		ParentID:  0,
		Type:      1,
	}
	{
		req := GetRequest(http.MethodPost, "/department/1/asset_class", headerJsonToken, GetJsonBody(CreateAssetClass))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	CreateAssetClass = define.CreateAssetClassReq{
		ClassName: "sub",
		ParentID:  1,
		Type:      1,
	}
	{
		req := GetRequest(http.MethodPost, "/department/1/asset_class", headerJsonToken, GetJsonBody(CreateAssetClass))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	// POST("/:department_id/asset
	CreateAsset2 := define.CreateAssetReq{
		AssetName:   "sub",
		Price:       decimal.New(123, 0),
		Description: "oo",
		Position:    "bb",
		ClassID:     1,
		Number:      1,
		Type:        1,
	}
	CreateAsset := define.CreateAssetReq{
		AssetName:   "name",
		Price:       decimal.New(123, 0),
		Description: "oo",
		Position:    "bb",
		ClassID:     1,
		Number:      1,
		Type:        1,
		Children:    []*define.CreateAssetReq{&CreateAsset2},
	}
	{
		req := GetRequest(http.MethodPost, "/department/1/asset", headerJsonToken, GetJsonBody(CreateAsset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, print_errormessage(res))
	}
	{
		req := GetRequest(http.MethodPost, "/department/1/asset", headerJson, GetJsonBody(CreateAsset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/department/1/asset", headerJsonToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/department/9/asset", headerJsonToken, GetJsonBody(CreateAsset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	// GET("/:department_id/asset/list
	{
		req := GetRequest(http.MethodGet, "/department/1/asset/list", headerJsonToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/department/1/asset/list", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/department/9/asset/list", headerJsonToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	// PATCH("/:department_id/asset/:asset_id"
	parent_id := uint(0)
	ModifyAsset := define.ModifyAssetInfoReq{
		AssetName:   "name",
		ParentID:    &parent_id,
		Price:       decimal.New(123, 0),
		Description: "oo",
		Position:    "bb",
		ClassID:     1,
		Number:      1,
		Type:        1,
	}
	{
		req := GetRequest(http.MethodPatch, "/department/1/asset/1", headerJsonToken, GetJsonBody(ModifyAsset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, print_errormessage(res))
	}
	{
		req := GetRequest(http.MethodPatch, "/department/1/asset/1", headerJson, GetJsonBody(ModifyAsset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPatch, "/department/1/asset/1", headerJsonToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPatch, "/department/9/asset/1", headerJsonToken, GetJsonBody(ModifyAsset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPatch, "/department/1/asset/9", headerJsonToken, GetJsonBody(ModifyAsset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	// PATCH("/:department_id/asset/expire
	expire := define.ExpireAssetReq{
		AssetID: 1,
	}
	{
		req := GetRequest(http.MethodPatch, "/department/1/asset/expire", headerJsonToken, GetJsonBody(expire))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, print_errormessage(res))
	}
	{
		req := GetRequest(http.MethodPatch, "/department/1/asset/expire", headerJson, GetJsonBody(expire))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPatch, "/department/1/asset/expire", headerJsonToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPatch, "/department/9/asset/expire", headerJsonToken, GetJsonBody(expire))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	// POST("/:department_id/asset/transfer"
	transfer := define.AssetTransferReq{
		UserID: 1,
		Assets: []define.ExpireAssetReq{expire},
	}
	{
		req := GetRequest(http.MethodPost, "/department/1/asset/transfer", headerJsonToken, GetJsonBody(transfer))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/department/1/asset/transfer", headerJson, GetJsonBody(transfer))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/department/1/asset/transfer", headerJsonToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/department/9/asset/transfer", headerJsonToken, GetJsonBody(transfer))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
}
