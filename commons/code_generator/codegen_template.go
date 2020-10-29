package code_generator

import "html/template"

var repositoryTmpl = template.Must(template.New("repository").Parse(`
package repositories

import (
	"github.com/jimersylee/iris-seed/commons"	
	"{{.PkgName}}/models"
	"github.com/jinzhu/gorm"
)

var {{.Name}}Repository = new{{.Name}}Repository()

func new{{.Name}}Repository() *{{.CamelName}}Repository {
	return &{{.CamelName}}Repository{}
}

type {{.CamelName}}Repository struct {
}

func (this *{{.CamelName}}Repository) Get(db *gorm.DB, id int64) *models.{{.Name}} {
	ret := &models.{{.Name}}{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *{{.CamelName}}Repository) Take(db *gorm.DB, where ...interface{}) *models.{{.Name}} {
	ret := &models.{{.Name}}{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *{{.CamelName}}Repository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.{{.Name}}) {
	cnd.Find(db, &list)
	return
}

func (this *{{.CamelName}}Repository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) ( ret *models.{{.Name}}) {
	cnd.FindOne(db, &ret)
	return 
}

func (this *{{.CamelName}}Repository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.{{.Name}}, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *{{.CamelName}}Repository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.{{.Name}}, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.{{.Name}}{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *{{.CamelName}}Repository) Create(db *gorm.DB, t *models.{{.Name}}) (err error) {
	err = db.Create(t).Error
	return
}

func (this *{{.CamelName}}Repository) Update(db *gorm.DB, t *models.{{.Name}}) (err error) {
	err = db.Save(t).Error
	return
}

func (this *{{.CamelName}}Repository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.{{.Name}}{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *{{.CamelName}}Repository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.{{.Name}}{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *{{.CamelName}}Repository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.{{.Name}}{}, "id = ?", id)
}

`))

var serviceTmpl = template.Must(template.New("service").Parse(`
package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"{{.PkgName}}/models"
	"{{.PkgName}}/repositories"
)

var {{.Name}}Service = new{{.Name}}Service()

func new{{.Name}}Service() *{{.CamelName}}Service {
	return &{{.CamelName}}Service {}
}

type {{.CamelName}}Service struct {
}

func (this *{{.CamelName}}Service) Get(id int64) *models.{{.Name}} {
	return repositories.{{.Name}}Repository.Get(db.GetDB(), id)
}

func (this *{{.CamelName}}Service) Take(where ...interface{}) *models.{{.Name}} {
	return repositories.{{.Name}}Repository.Take(db.GetDB(), where...)
}

func (this *{{.CamelName}}Service) Find(cnd *commons.SqlCnd) []models.{{.Name}} {
	return repositories.{{.Name}}Repository.Find(db.GetDB(), cnd)
}

func (this *{{.CamelName}}Service) FindOne(cnd *commons.SqlCnd) *models.{{.Name}} {
	return repositories.{{.Name}}Repository.FindOne(db.GetDB(), cnd)
}

func (this *{{.CamelName}}Service) FindPageByParams(params *commons.QueryParams) (list []models.{{.Name}}, paging *commons.Paging) {
	return repositories.{{.Name}}Repository.FindPageByParams(db.GetDB(), params)
}

func (this *{{.CamelName}}Service) FindPageByCnd(cnd *commons.SqlCnd) (list []models.{{.Name}}, paging *commons.Paging) {
	return repositories.{{.Name}}Repository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *{{.CamelName}}Service) Create(t *models.{{.Name}}) error {
	return repositories.{{.Name}}Repository.Create(db.GetDB(), t)
}

func (this *{{.CamelName}}Service) Update(t *models.{{.Name}}) error {
	return repositories.{{.Name}}Repository.Update(db.GetDB(), t)
}

func (this *{{.CamelName}}Service) Updates(id int64, columns map[string]interface{}) error {
	return repositories.{{.Name}}Repository.Updates(db.GetDB(), id, columns)
}

func (this *{{.CamelName}}Service) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.{{.Name}}Repository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *{{.CamelName}}Service) Delete(id int64) {
	repositories.{{.Name}}Repository.Delete(db.GetDB(), id)
}

`))

var controllerTmpl = template.Must(template.New("controller").Parse(`
package api

import (
	"{{.PkgName}}/models"
	"{{.PkgName}}/services"
	"github.com/jimersylee/iris-seed/commons"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
)

type {{.Name}}Controller struct {
	Ctx             iris.Context
}

func (this *{{.Name}}Controller) GetBy(id int64) *response.WebApiRes {
	t := services.{{.Name}}Service.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}
	return response.JsonData(t)
}

func (this *{{.Name}}Controller) AnyList() *response.WebApiRes {
	list, paging := services.{{.Name}}Service.FindPageByParams(commons.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return response.JsonData(&commons.PageResult{Results: list, Page: paging})
}

func (this *{{.Name}}Controller) PostCreate() *response.WebApiRes {
	t := &models.{{.Name}}{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.{{.Name}}Service.Create(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}

func (this *{{.Name}}Controller) PostUpdate() *response.WebApiRes {
	id, err := commons.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	t := services.{{.Name}}Service.Get(id)
	if t == nil {
			return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.{{.Name}}Service.Update(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}

`))
