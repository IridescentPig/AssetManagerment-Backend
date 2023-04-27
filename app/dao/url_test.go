package dao

import (
	"asset-management/app/model"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrl(t *testing.T) {
	Init()
	entity := model.Entity{
		Name: "test_entity",
	}
	err := EntityDao.Create(entity)
	assert.Equal(t, nil, err, "database error")

	url := model.Url{
		Name:            "0",
		Url:             "www.com",
		EntityID:        1,
		DepartmentSuper: false,
		EntitySuper:     false,
		SystemSuper:     false,
	}
	UrlDao.Create(url)

	url = model.Url{
		Name:            "1",
		Url:             "www.com",
		EntityID:        1,
		DepartmentSuper: false,
		EntitySuper:     false,
		SystemSuper:     true,
	}
	UrlDao.Create(url)

	url = model.Url{
		Name:            "2",
		Url:             "www.com",
		EntityID:        1,
		DepartmentSuper: false,
		EntitySuper:     true,
		SystemSuper:     false,
	}
	UrlDao.Create(url)

	url = model.Url{
		Name:            "3",
		Url:             "www.com",
		EntityID:        1,
		DepartmentSuper: false,
		EntitySuper:     true,
		SystemSuper:     true,
	}
	UrlDao.Create(url)

	url = model.Url{
		Name:            "4",
		Url:             "www.com",
		EntityID:        1,
		DepartmentSuper: true,
		EntitySuper:     false,
		SystemSuper:     false,
	}
	UrlDao.Create(url)

	url = model.Url{
		Name:            "5",
		Url:             "www.com",
		EntityID:        1,
		DepartmentSuper: true,
		EntitySuper:     false,
		SystemSuper:     true,
	}
	UrlDao.Create(url)

	url = model.Url{
		Name:            "6",
		Url:             "www.com",
		EntityID:        1,
		DepartmentSuper: true,
		EntitySuper:     true,
		SystemSuper:     false,
	}
	UrlDao.Create(url)

	url = model.Url{
		Name:            "7",
		Url:             "www.com",
		EntityID:        1,
		DepartmentSuper: true,
		EntitySuper:     true,
		SystemSuper:     true,
	}
	UrlDao.Create(url)

	assets, err := UrlDao.GetUrlsByEntity(1, false, false, false)
	assert.Equal(t, nil, err, "database error")
	log.Print("------")
	for _, asset := range assets {
		log.Print(asset.Name)
	}
	assets, err = UrlDao.GetUrlsByEntity(1, false, false, true)
	assert.Equal(t, nil, err, "database error")
	log.Print("------")
	for _, asset := range assets {
		log.Print(asset.Name)
	}
	assets, err = UrlDao.GetUrlsByEntity(1, false, true, false)
	assert.Equal(t, nil, err, "database error")
	log.Print("------")
	for _, asset := range assets {
		log.Print(asset.Name)
	}
	assets, err = UrlDao.GetUrlsByEntity(1, false, true, true)
	assert.Equal(t, nil, err, "database error")
	log.Print("------")
	for _, asset := range assets {
		log.Print(asset.Name)
	}
	assets, err = UrlDao.GetUrlsByEntity(1, true, false, false)
	assert.Equal(t, nil, err, "database error")
	log.Print("------")
	for _, asset := range assets {
		log.Print(asset.Name)
	}
	assets, err = UrlDao.GetUrlsByEntity(1, true, false, true)
	assert.Equal(t, nil, err, "database error")
	log.Print("------")
	for _, asset := range assets {
		log.Print(asset.Name)
	}
	assets, err = UrlDao.GetUrlsByEntity(1, true, true, false)
	assert.Equal(t, nil, err, "database error")
	log.Print("------")
	for _, asset := range assets {
		log.Print(asset.Name)
	}
	assets, err = UrlDao.GetUrlsByEntity(1, true, true, true)
	assert.Equal(t, nil, err, "database error")
	log.Print("------")
	for _, asset := range assets {
		log.Print(asset.Name)
	}
}
