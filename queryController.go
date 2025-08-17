/***************************************************
 ** @Desc : This file for ...
 ** @Time : 2019/11/6 14:03
 ** @Author : yuebin
 ** @File : query.go
 ** @Last Modified by : yuebin
 ** @Last Modified time: 2019/11/6 14:03
 ** @Software: GoLand
****************************************************/
package controllers

import (
	"boss/service"
	"github.com/beego/beego/v2/server/web"
	"strings"
)

type SupplierQuery struct {
	web.Controller
}

func (c *SupplierQuery) SupplierOrderQuery() {

	bankOrderId := strings.TrimSpace(c.GetString("bankOrderId"))

	se := new(service.QueryService)
	keyDataJSON := se.SupplierOrderQuery(bankOrderId)

	c.Data["json"] = keyDataJSON
	_ = c.ServeJSON()
}

/*
* 向上游查询代付结果
 */
func (c *SupplierQuery) SupplierPayForQuery() {
	bankOrderId := strings.TrimSpace(c.GetString("bankOrderId"))

	se := new(service.QueryService)
	keyDataJSON := se.SupplierPayForQuery(bankOrderId)

	c.Data["json"] = keyDataJSON
	_ = c.ServeJSON()
}
