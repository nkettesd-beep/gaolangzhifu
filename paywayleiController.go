/***************************************************
 ** @Desc : This file for ...
 ** @Time : 2019/10/29 15:01
 ** @Author : yuebin
 ** @File : pay_way_code
 ** @Last Modified by : yuebin
 ** @Last Modified time: 2019/10/29 15:01
 ** @Software: GoLand
****************************************************/
package controllers

import (
	"boss/datas"
	"boss/models/accounts"
	"boss/models/road"
	"boss/utils"
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/rs/xid"
	"github.com/widuu/gojson"
	"strconv"
	"strings"
	"time"
)

type PaywayleiController struct {
	BaseController
	DisplayCount int
	CurrentPage  int
	TotalPage    int
	JumpPage     int
	Offset       int
}



func (c *PaywayleiController)GetScanPayWayCodes()  {
	menuDataJSON := new(datas.ProDataJSON)
	se := new(road.RoadBandlei)
	roadbank:=se.Getroadband()

	if len(roadbank) == 0 {
		menuDataJSON.Code = -1
	} else {
		menuDataJSON.Code = 200
	}
	menuDataJSON.Roadbank = roadbank

	c.GenerateJSON(menuDataJSON)

}

func (c *PaywayleiController)GetScanPayWayCodesinfo()  {
	menuDataJSON := new(datas.ProDataJSON)
	se := new(road.RoadInfo)
	roadbank:=se.Getroadbandinfo()

	if len(roadbank) == 0 {
		menuDataJSON.Code = -1
	} else {
		menuDataJSON.Code = 200
	}
	menuDataJSON.Roadbankinfo = roadbank

	c.GenerateJSON(menuDataJSON)

}


func (c *PaywayleiController) GetCutPage(l int) {
	c.DisplayCount, _ = c.GetInt("displayCount")
	c.CurrentPage, _ = c.GetInt("currentPage")
	c.TotalPage, _ = c.GetInt("totalPage")
	c.JumpPage, _ = c.GetInt("jumpPage")

	if c.CurrentPage == 0 {
		c.CurrentPage = 1
	}
	if c.DisplayCount == 0 {
		c.DisplayCount = 20
	}
	if c.JumpPage > 0 {
		c.CurrentPage = c.JumpPage
	}

	if l > 0 {
		c.TotalPage = l / c.DisplayCount
		if l%c.DisplayCount > 0 {
			c.TotalPage += 1
		}
	} else {
		c.TotalPage = 0
		c.CurrentPage = 0
	}
	//假如当前页超过了总页数
	if c.CurrentPage > c.TotalPage {
		c.CurrentPage = c.TotalPage
	}
	//计算出偏移量
	c.Offset = (c.CurrentPage - 1) * c.DisplayCount
}

func (c *PaywayleiController) GetAccounttody() {
	accountName := strings.TrimSpace(c.GetString("accountName"))
	accountUid := strings.TrimSpace(c.GetString("accountNo"))
	start := strings.TrimSpace(c.GetString("startTime"))
	end := strings.TrimSpace(c.GetString("endTime"))
	typecode := strings.TrimSpace(c.GetString("tongdao"))
	params := make(map[string]string)
	params["account_name__icontains"] = accountName
	params["account_uid_icontains"] = accountUid

	l := Accounts.GetAccountLenByMap(params)
	c.GetCutPage(l)

	accountDataJSON := new(datas.AccountDataJSON)
	accountDataJSON.DisplayCount = c.DisplayCount
	accountDataJSON.Code = 200
	accountDataJSON.CurrentPage = c.CurrentPage
	accountDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		accountDataJSON.AccountList = make([]Accounts.AccountInfo, 0)
		c.GenerateJSON(accountDataJSON)
		return
	}

	accountDataJSON.StartIndex = c.Offset
	accountDataJSON.AccountTody = Accounts.GetAccountByMappro(params, c.DisplayCount, c.Offset,start,end,accountName,typecode)

	accountDataJSON.AccountCount = Accounts.GetAccountByMapprocount(params, c.DisplayCount, c.Offset,start,end,accountName,typecode)
	accountDataJSON.AccountCountsuccess = Accounts.GetAccountByMapprocountsuccess(params, c.DisplayCount, c.Offset,start,end,accountName,typecode)
	c.GenerateJSON(accountDataJSON)
}

type Pamelei struct {
	Params string
	TelegramBankamount float64

}
type ScanDatalei struct {
	Supplier   string //上游的通道供应商
	PayType    string //支付类型
	OrderNo    string //下游商户请求订单号
	BankNo     string //本系统的请求订单号
	OrderPrice string //订单金额
	FactPrice  string //实际的展示在客户面前的金额
	Status     string //状态码 '00' 成功
	PayUrl     string //支付二维码链接地址
	Msg        string //附加的信息
}
func (c *PaywayleiController) GetAccounttodybug() {
	keyDataJSONv := new(datas.RoadPoolDataJSON)
	var u Pame
    json.Unmarshal(c.Ctx.Input.RequestBody,&u)


	mchNo := gojson.Json(u.Params).Get("mchNo").Tostring()
	apiKey := gojson.Json(u.Params).Get("apikey").Tostring()
	signKey := gojson.Json(u.Params).Get("signkey").Tostring()
	bankcode := gojson.Json(u.Params).Get("bankcode").Tostring()
	controller := gojson.Json(u.Params).Get("controller").Tostring()
	clientip := gojson.Json(u.Params).Get("clientip").Tostring()   //扩展
	device := gojson.Json(u.Params).Get("device").Tostring()    //扩展
	param1 := gojson.Json(u.Params).Get("param1").Tostring()      //扩展
	param2 := gojson.Json(u.Params).Get("param2").Tostring()      //扩展
	paramkey := gojson.Json(u.Params).Get("paramkey").Tostring()      //扩展
	paramkey2 := gojson.Json(u.Params).Get("paramkey2").Tostring()      //扩展
	amount :=u.TelegramBankamount

	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Millisecond)
	params := make(map[string]string)
	params["mchNo"] = mchNo
	params["amount"] = strconv.FormatInt(int64(amount*100), 10)
	params["reqTime"] =strconv.FormatInt(timestamp,10)
	params["apikey"] = apiKey
	params["bankcode"] = bankcode
	params["mchOrderNo"] = strconv.FormatInt(timestamp,10)
	params["signkey"] =signKey
	params["notifyUrl"] =mchNo +"/"+controller+".php"
	params["controller"] =controller
	params["clientip"] =clientip
	params["device"] =device
	params["param1"] =param1
	params["param2"] =param2
	params["paramkey"] =paramkey
	params["paramkey2"] =paramkey2


	request := mchNo +"/"+controller+".php?" + utils.MapToString(params)

	logs.Info("请求字符串 = " + request)

	var scanData ScanData
	//scanData.Status = "00"
	response, err := httplib.Post(request).String()


	//order.Bankorderhuitiao(orderInfo.BankOrderId, gojson.Json(response).Get("msg").Tostring())
	if err != nil {
		logs.Error("支付请求失败：" + err.Error())
		scanData.Status = "-1"
		scanData.Msg = "请求失败：" + err.Error()
	} else {
		logs.Info("支付返回 = " + response)
		//status := gojson.Json(response).Get("code").Tostring()
		//message := gojson.Json(response).Get("msg").Tostring()
		//codeUrl := gojson.Json(response).Get("payurl").Tostring()
		//fmt.Print(codeUrl)
		//if "0" != status {
		//	scanData.Status = "-1"
		//	scanData.Msg = message
		//} else {
		//
		//	//fmt.Print(codeUrl)
		//	//scanData.PayUrl = codeUrl
		//	//scanData.OrderNo = orderInfo.BankOrderId
		//	//scanData.OrderPrice = strconv.FormatFloat(orderInfo.OrderAmount, 'f', 2, 64)
		//}
	}
	//fmt.Print(scanData)
	//return scanData
	keyDataJSONv.Code = 200
	keyDataJSONv.Msg= response


	//accountDataJSON.StartIndex = c.Offset
	//accountDataJSON.AccountTody = Accounts.GetAccountByMappro(params, c.DisplayCount, c.Offset,start,end,accountName,typecode)
	c.GenerateJSON(keyDataJSONv)
}

/*
* 添加操作员
 */
func (c *PaywayleiController) AddOperator() {
	roadName := "lei" + xid.New().String()
	name := strings.TrimSpace(c.GetString("RoadName"))
	//role := strings.TrimSpace(c.GetString("operatorRole"))
	//status := strings.TrimSpace(c.GetString("status"))
	//remark := strings.TrimSpace(c.GetString("remark"))
	keyDataJSONv := new(datas.RoadPoolDataJSON)
	if roadName==""{
		keyDataJSONv.Code = -1
		c.GenerateJSON(keyDataJSONv)
	}

	se := new(road.RoadBandlei)
	keyDataJSON := se.AddOperator(roadName,name)

	if keyDataJSON==true{
		keyDataJSONv.Code = 200
		keyDataJSONv.Msg ="添加成功"

		c.GenerateJSON(keyDataJSONv)
	}

}


func (c *PaywayleiController) DeleteMenu() {
	menuUid := c.GetString("roadPoolCode")
	se := new(road.RoadBandlei)
	keyDataJSONv := new(datas.RoadPoolDataJSON)
	dataJSON := se.DeleteMenuInfo(menuUid)

	if dataJSON==true{
		keyDataJSONv.Code = 200
		keyDataJSONv.Msg ="删除成功"

		c.GenerateJSON(keyDataJSONv)
	}
}

