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
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/rs/xid"
	"github.com/widuu/gojson"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type PaywayController struct {
	BaseController
	DisplayCount int
	CurrentPage  int
	TotalPage    int
	JumpPage     int
	Offset       int
}



//var ScanPayWayCodes = []string{
//	"WEIXIN_SCAN",
//	"UNION_SCAN",
//	"ALI_SCAN",
//	"BAIDU_SCAN",
//	"JD_SCAN",
//	"QQ_SCAN",
//	"Aw_SCAN",
//}

var H5PayWayCodes = []string{
	"WEIXIN_H5",
	"ALI_H5",
	"QQ_H5",
	"UNION_H5",
	"BAIDU_H5",
	"JD_H5",
}

var SytPayWayCodes = []string{
	"WEIXIN_SYT",
	"ALI_SYT",
	"QQ_SYT",
	"UNION_SYT",
	"BAIDU_SYT",
	"JD_SYT",
}

var FastPayWayCodes = []string{
	"UNION-FAST",
}



var WebPayWayCode = []string{
	"UNION-WAP",
}

func (c *PaywayController)GetScanPayWayCodes()  {
	menuDataJSON := new(datas.ProDataJSON)
	se := new(road.RoadBand)
	roadbank:=se.Getroadband()

	if len(roadbank) == 0 {
		menuDataJSON.Code = -1
	} else {
		menuDataJSON.Code = 200
	}
	menuDataJSON.Roadbank = roadbank

	c.GenerateJSON(menuDataJSON)

}

func (c *PaywayController)GetScanPayWayCodesinfo()  {
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


func (c *PaywayController) GetCutPage(l int) {
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

func (c *PaywayController) GetAccounttody() {
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

type Pame struct {
	Params string
	TelegramBankamount float64
	PayType string
	ProductNamecode string
	Userid string


}
type ScanData struct {
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
func (c *PaywayController) GetAccounttodybug() {
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


func (c *PaywayController)TestPaycheshi() {
	keyDataJSONv := new(datas.RoadPoolDataJSON)
	var uv Pame
	json.Unmarshal(c.Ctx.Input.RequestBody,&uv)


	Userid := uv.Userid
	ProductNamecode := uv.ProductNamecode
	orderPrice :=uv.TelegramBankamount
	PayType := uv.PayType
	fmt.Print(PayType)
	params := make(map[string]string)
	params["orderNo"] = xid.New().String()
	params["productName"] = "kongyuhebin"
	params["orderPeriod"] = "1"
	params["orderPrice"] = strconv.FormatInt(int64(orderPrice), 10)
	params["payWayCode"] = PayType
	params["osType"] = "1"
	params["notifyUrl"] = "http://localhost:12309/shop/notify"
	params["payKey"] =Userid
	keys := utils.SortMap(params)
	params["sign"] = utils.GetMD5Sign(params, keys, ProductNamecode)

	u := url.Values{}
	for k, v := range params {
		u.Add(k, v)
	}

	l := "http://localhost:12309/gateway/scan?" + u.Encode()
	logs.Info("请求url：" + l)

	resp := httplib.Get(l)
	s, err := resp.String()

	if err != nil {
		logs.Error("请求错误：" + err.Error())

	}

	logs.Info("微信扫码返回结果：" + s)
	keyDataJSONv.Code = 200
	keyDataJSONv.Msg= s


	//accountDataJSON.StartIndex = c.Offset
	//accountDataJSON.AccountTody = Accounts.GetAccountByMappro(params, c.DisplayCount, c.Offset,start,end,accountName,typecode)
	c.GenerateJSON(keyDataJSONv)
}

/*
* 添加操作员
 */
func (c *PaywayController) AddOperator() {
	roadName := strings.TrimSpace(c.GetString("RoadRemark"))
	name := strings.TrimSpace(c.GetString("RoadName"))
	//role := strings.TrimSpace(c.GetString("operatorRole"))
	//status := strings.TrimSpace(c.GetString("status"))
	//remark := strings.TrimSpace(c.GetString("remark"))
	keyDataJSONv := new(datas.RoadPoolDataJSON)
	if roadName==""{
		keyDataJSONv.Code = -1
		c.GenerateJSON(keyDataJSONv)
	}

	se := new(road.RoadBand)
	keyDataJSON := se.AddOperator(roadName,name)

	if keyDataJSON==true{
		keyDataJSONv.Code = 200
		keyDataJSONv.Msg ="添加成功"

		c.GenerateJSON(keyDataJSONv)
	}

}


func (c *PaywayController) DeleteMenu() {
	menuUid := c.GetString("roadPoolCode")
	se := new(road.RoadBand)
	keyDataJSONv := new(datas.RoadPoolDataJSON)
	dataJSON := se.DeleteMenuInfo(menuUid)

	if dataJSON==true{
		keyDataJSONv.Code = 200
		keyDataJSONv.Msg ="删除成功"

		c.GenerateJSON(keyDataJSONv)
	}
}

func GetNameByPayWayCode(code string) string {
	switch code {
	case "WEIXIN_SCAN":
		return "微信扫码"
	case "UNION_SCAN":
		return "银联扫码"
	case "ALI_SCAN":
		return "支付宝扫码"
	case "BAIDU_SCAN":
		return "百度扫码"
	case "JD_SCAN":
		return "京东扫码"
	case "QQ_SCAN":
		return "QQ扫码"

	case "WEIXIN_H5":
		return "微信H5"
	case "UNION_H5":
		return "银联H5"
	case "ALI_H5":
		return "支付宝H5"
	case "BAIDU_H5":
		return "百度H5"
	case "JD_H5":
		return "京东H5"
	case "QQ_H5":
		return "QQ-H5"

	case "WEIXIN_SYT":
		return "微信收银台"
	case "UNION_SYT":
		return "银联收银台"
	case "ALI_SYT":
		return "支付宝收银台"
	case "BAIDU_SYT":
		return "百度收银台"
	case "JD_SYT":
		return "京东收银台"
	case "QQ_SYT":
		return "QQ-收银台"

	case "UNION_FAST":
		return "银联快捷"
	case "UNION_WAP":
		return "银联web"
	default:
		return "未知"
	}
}
