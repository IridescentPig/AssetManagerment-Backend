package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"

	"github.com/gin-gonic/gin/binding"
)

type urlApi struct {
}

var UrlApi *urlApi

func newUrlApi() *urlApi {
	return &urlApi{}
}

func init() {
	UrlApi = newUrlApi()
}

/*
Handle func for POST /entity/{entity_id}/url
*/
func (url *urlApi) CreateUrl(ctx *utils.Context) {
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}

	exists, err := service.EntityService.ExistsEntityByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return
	}

	entitySuper := service.UserService.EntitySuper(ctx)
	if !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	var req define.CreateUrlReq
	err = ctx.MustBindWith(&req, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	exists, err = service.UrlService.CheckIfUrlExists(req.Name, entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if exists {
		ctx.BadRequest(myerror.DUPLICATED_NAME, myerror.DUPLICATED_NAME_INFO)
		return
	}

	err = service.UrlService.CreateUrl(req, entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	ctx.Success(nil)
}

/*
Handle func for PATCH /entity/{entity_id}/url
*/
func (url *urlApi) ModifyUrl(ctx *utils.Context) {
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}

	exists, err := service.EntityService.ExistsEntityByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return
	}

	entitySuper := service.UserService.EntitySuper(ctx)
	if !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	var req define.ModifyUrlReq
	err = ctx.MustBindWith(&req, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	exists, err = service.UrlService.CheckIfUrlExists(req.OldName, entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(myerror.URL_NOT_FOUND, myerror.URL_NOT_FOUND_INFO)
		return
	}

	if req.Name != req.OldName {
		exists, err = service.UrlService.CheckIfUrlExists(req.Name, entityID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		} else if exists {
			ctx.BadRequest(myerror.DUPLICATED_NAME, myerror.DUPLICATED_NAME_INFO)
			return
		}
	}

	err = service.UrlService.ModifyUrlInfo(req, entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	ctx.Success(nil)
}

/*
Handle func for DELETE /entity/{entity_id}/url
*/
func (url *urlApi) DeleteUrl(ctx *utils.Context) {
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}

	exists, err := service.EntityService.ExistsEntityByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return
	}

	entitySuper := service.UserService.EntitySuper(ctx)
	if !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	var req define.DeleteUrlReq
	err = ctx.MustBindWith(&req, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	exists, err = service.UrlService.CheckIfUrlExists(req.Name, entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(myerror.URL_NOT_FOUND, myerror.URL_NOT_FOUND_INFO)
		return
	}

	err = service.UrlService.DeleteUrl(req.Name, entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	ctx.Success(nil)
}

/*
Handle func for GET /entity/{entity_id}/url
*/
func (url *urlApi) GetUrlsByEntity(ctx *utils.Context) {
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}

	exists, err := service.EntityService.ExistsEntityByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return
	}

	departmentSuper := service.UserService.DepartmentSuper(ctx)
	entitySuper := service.UserService.EntitySuper(ctx)
	systemSuper := service.UserService.SystemSuper(ctx)
	url_list, err := service.UrlService.GetUrlsByEntity(entityID, departmentSuper, entitySuper, systemSuper)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	var response define.GetUrlResponse
	for _, target_url := range url_list {
		response.UrlList = append(response.UrlList, define.UrlInfo{
			Name:            target_url.Name,
			Url:             target_url.Url,
			DepartmentSuper: target_url.DepartmentSuper,
			EntitySuper:     target_url.EntitySuper,
			SystemSuper:     target_url.SystemSuper,
		})
	}

	ctx.Success(response)
}

/*
Handle func for GET /entity/{entity_id}/url/list
*/
func (url *urlApi) GetUrlList(ctx *utils.Context) {
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}

	exists, err := service.EntityService.ExistsEntityByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return
	}

	entitySuper := service.UserService.EntitySuper(ctx)
	if !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	url_list, err := service.UrlService.GetUrlsByEntity(entityID, true, true, true)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	var response define.GetUrlResponse
	for _, target_url := range url_list {
		response.UrlList = append(response.UrlList, define.UrlInfo{
			Name:            target_url.Name,
			Url:             target_url.Url,
			DepartmentSuper: target_url.DepartmentSuper,
			EntitySuper:     target_url.EntitySuper,
			SystemSuper:     target_url.SystemSuper,
		})
	}

	ctx.Success(response)
}
