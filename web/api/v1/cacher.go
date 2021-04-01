package v1

import (
	"net/http"
	"time"

	"github.com/powerpuffpenguin/xormcache/cache"
	"github.com/powerpuffpenguin/xormcache/web"

	"github.com/gin-gonic/gin"
)

// Cacher .
type Cacher struct {
	web.Helper
}

// Register impl IHelper
func (h Cacher) Register(router *gin.RouterGroup) {
	r := router.Group(`cacher`)
	r.GET(`ids`, h.getIds)
	r.GET(`bean`, h.getBean)
	r.PUT(`ids`, h.putIds)
	r.PUT(`bean`, h.putBean)
	r.DELETE(`ids`, h.delIds)
	r.DELETE(`bean`, h.delBean)
	r.DELETE(`ids/:tableName`, h.clearIds)
	r.DELETE(`beans/:tableName`, h.clearBeans)
	r.GET(`detail`, h.detail)
}
func (h Cacher) getIds(c *gin.Context) {
	var obj struct {
		TableName string `form:"tableName" json:"tableName" xml:"tableName" yaml:"tableName"`
		SQL       string `form:"sql" json:"sql" xml:"sql" yaml:"sql"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	cacher := cache.DefaultCacher()
	val := cacher.GetIds(obj.TableName, obj.SQL)
	if val == nil {
		c.Status(http.StatusNotFound)
	} else {
		ele := val.(cache.Element)
		h.NegotiateObject(c, ele.Modtime, gin.H{
			`data`: ele.Data,
		})
	}
}
func (h Cacher) getBean(c *gin.Context) {
	var obj struct {
		TableName string `form:"tableName" json:"tableName" xml:"tableName" yaml:"tableName"`
		ID        string `form:"id" json:"id" xml:"id" yaml:"id"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	cacher := cache.DefaultCacher()
	val := cacher.GetBean(obj.TableName, obj.ID)
	if val == nil {
		c.Status(http.StatusNotFound)
	} else {
		ele := val.(cache.Element)
		h.NegotiateObject(c, ele.Modtime, gin.H{
			`data`: ele.Data,
		})
	}
}
func (h Cacher) putIds(c *gin.Context) {
	var obj struct {
		TableName string `form:"tableName" json:"tableName" xml:"tableName" yaml:"tableName"`
		SQL       string `form:"sql" json:"sql" xml:"sql" yaml:"sql"`
		IDS       []byte `form:"ids" json:"ids" xml:"ids" yaml:"ids"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	cacher := cache.DefaultCacher()
	cacher.PutIds(obj.TableName, obj.SQL, obj.IDS)
}
func (h Cacher) putBean(c *gin.Context) {
	var obj struct {
		TableName string `form:"tableName" json:"tableName" xml:"tableName" yaml:"tableName"`
		ID        string `form:"id" json:"id" xml:"id" yaml:"id"`
		Obj       []byte `form:"obj" json:"obj" xml:"obj" yaml:"obj"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	cacher := cache.DefaultCacher()
	cacher.PutBean(obj.TableName, obj.ID, obj.Obj)
}
func (h Cacher) delIds(c *gin.Context) {
	var obj struct {
		TableName string `form:"tableName" json:"tableName" xml:"tableName" yaml:"tableName"`
		SQL       string `form:"sql" json:"sql" xml:"sql" yaml:"sql"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	cacher := cache.DefaultCacher()
	cacher.DelIds(obj.TableName, obj.SQL)
}
func (h Cacher) delBean(c *gin.Context) {
	var obj struct {
		TableName string `form:"tableName" json:"tableName" xml:"tableName" yaml:"tableName"`
		ID        string `form:"id" json:"id" xml:"id" yaml:"id"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	cacher := cache.DefaultCacher()
	cacher.DelBean(obj.TableName, obj.ID)
}
func (h Cacher) clearIds(c *gin.Context) {
	var obj struct {
		TableName string `uri:"tableName"`
	}
	e := h.BindURI(c, &obj)
	if e != nil {
		return
	}
	cacher := cache.DefaultCacher()
	cacher.ClearIds(obj.TableName)
}
func (h Cacher) clearBeans(c *gin.Context) {
	var obj struct {
		TableName string `uri:"tableName"`
	}
	e := h.BindURI(c, &obj)
	if e != nil {
		return
	}
	cacher := cache.DefaultCacher()
	cacher.ClearBeans(obj.TableName)
}
func (h Cacher) detail(c *gin.Context) {
	cacher := cache.DefaultCacher()
	h.NegotiateObject(c, startAt, gin.H{
		`maxAge`:         uint32(cacher.Expired / time.Second),
		`maxElementSize`: uint32(cacher.MaxElementSize),
	})
}
