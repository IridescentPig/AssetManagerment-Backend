package api

import "asset-management/utils"

type urlApi struct {
}

var UrlApi *urlApi

func newUrlApi() *urlApi {
	return &urlApi{}
}

func init() {
	UrlApi = newUrlApi()
}

func (url *urlApi) CreateUrl(ctx *utils.Context) {

}

func (url *urlApi) ModifyUrl(ctx *utils.Context) {

}

func (url *urlApi) DeleteUrl(ctx *utils.Context) {

}

func (url *urlApi) GetUrlsByEntity(ctx *utils.Context) {

}
