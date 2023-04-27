package service

import (
	"asset-management/app/define"
	"testing"
)

func TestUrl(t *testing.T) {
	InitForTest()

	_ = EntityService.CreateEntity("test_entity")

	req := define.CreateUrlReq{
		Name:            "0",
		Url:             "www.com",
		DepartmentSuper: false,
		EntitySuper:     false,
		SystemSuper:     false,
	}
	UrlService.CreateUrl(req, 1)

	req = define.CreateUrlReq{
		Name:            "2",
		Url:             "www.com",
		DepartmentSuper: true,
		EntitySuper:     false,
		SystemSuper:     false,
	}
	UrlService.CreateUrl(req, 1)

	modify_req := define.ModifyUrlReq{
		OldName:         "0",
		Name:            "0",
		Url:             "wwwwwww.com",
		DepartmentSuper: false,
		EntitySuper:     false,
		SystemSuper:     false,
	}
	UrlService.ModifyUrlInfo(modify_req, 1)

	modify_req = define.ModifyUrlReq{
		OldName:         "1",
		Name:            "qqqq",
		Url:             "wwwwwww.com",
		DepartmentSuper: true,
		EntitySuper:     false,
		SystemSuper:     false,
	}
	UrlService.ModifyUrlInfo(modify_req, 1)

	UrlService.GetUrlsByEntity(1, true, true, true)
	UrlService.GetUrlsByEntity(1, false, false, true)

	UrlService.CheckIfUrlExists("0", 1)
	UrlService.CheckIfUrlExists("qqq", 1)

	UrlService.DeleteUrl("unexist", 1)
	UrlService.DeleteUrl("0", 1)
}
