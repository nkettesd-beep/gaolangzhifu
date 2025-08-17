package controllers

import (
	"boss/common"
	"boss/datas"
	"boss/models/accounts"
	"boss/models/agent"
	"boss/models/mauser"
	"boss/models/merchant"
	"boss/models/notify"
	"boss/models/order"
	"boss/models/payfor"
	"boss/models/road"
	"boss/models/roadband"
	"boss/models/system"
	"boss/models/user"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TelegramController struct {
	BaseController
	DisplayCount int
	CurrentPage  int
	TotalPage    int
	JumpPage     int
	Offset       int
}

////跨域
//func (c *GetController) AllowCross() {
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")       //允许访问源
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, PUT, POST")    //允许post访问
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
//	c.Ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
//}


/*
* 处理分页的函数
 */
func (c *TelegramController) GettelegramCutPage(l int) {
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

func (c *TelegramController) GettelegramUid(){
	menuDataJSON := new(datas.MenuDataJSONjj)
	menuList := roadband.Getroadband()


	menuDataJSON.Code = 200
	menuDataJSON.MenuList = menuList
	c.GenerateJSON(menuDataJSON)
}

func (c *TelegramController) GettelegramMenu() {

	firstMenuSearch := strings.TrimSpace(c.GetString("firstMenuSearch"))

	params := make(map[string]string)
	params["first_menu__icontains"] = firstMenuSearch


	menuDataJSON := new(datas.MenuDataJSON)
	menuDataJSON.DisplayCount = c.DisplayCount
	menuDataJSON.CurrentPage = c.CurrentPage
	menuDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		menuDataJSON.Code = -1
		menuDataJSON.MenuList = make([]system.MenuInfo, 0)
		c.GenerateJSON(menuDataJSON)
		return
	}

	menuInfoList := system.GetMenuOffsetByMap(params, c.DisplayCount, c.Offset)
	sort.Sort(system.MenuInfoSlice(menuInfoList))
	for i, m := range menuInfoList {
		secondMenuInfoList := system.GetSecondMenuListByFirstMenuUid(m.MenuUid)
		smenus := ""
		for j := 0; j < len(secondMenuInfoList); j++ {
			smenus += secondMenuInfoList[j].SecondMenu
			if j != (len(secondMenuInfoList) - 1) {
				smenus += "||"
			}
		}
		m.SecondMenu = smenus
		menuInfoList[i] = m
	}
	menuDataJSON.Code = 200
	menuDataJSON.MenuList = menuInfoList
	menuDataJSON.StartIndex = c.Offset

	if len(menuInfoList) == 0 {
		menuDataJSON.Msg = "获取菜单列表失败"
	}

	c.GenerateJSON(menuDataJSON)
}

func (c *TelegramController) GettelegramFirstMenu() {
	menuDataJSON := new(datas.MenuDataJSON)
	menuList := system.GetMenuAll()

	if len(menuList) == 0 {
		menuDataJSON.Code = -1
	} else {
		menuDataJSON.Code = 200
	}
	sort.Sort(system.MenuInfoSlice(menuList))
	menuDataJSON.MenuList = menuList
	c.GenerateJSON(menuDataJSON)
}


func (c *TelegramController) Gettelegramcrollerd() {
	menuDataJSON := new(datas.MenuDataJSONjj)
	menuList := system.GetMenuAllorcer()
	menuDataJSON.Code = 200
	menuDataJSON.MenuList = menuList
	c.GenerateJSON(menuDataJSON)
}

/*
*获取所有的二级菜单
 */
func (c *TelegramController) GettelegramSecondMenu() {

	firstMenuSearch := strings.TrimSpace(c.GetString("firstMenuSerach"))
	secondMenuSearch := strings.TrimSpace(c.GetString("secondMenuSerach"))

	params := make(map[string]string)
	params["first_menu__icontains"] = firstMenuSearch
	params["second_menu__icontains"] = secondMenuSearch


	secondMenuDataJSON := new(datas.SecondMenuDataJSON)
	secondMenuDataJSON.DisplayCount = c.DisplayCount

	secondMenuDataJSON.Code = 200
	secondMenuDataJSON.CurrentPage = c.CurrentPage
	secondMenuDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		secondMenuDataJSON.SecondMenuList = make([]system.SecondMenuInfo, 0)
		c.GenerateJSON(secondMenuDataJSON)
		return
	}

	secondMenuList := system.GetSecondMenuByMap(params, c.DisplayCount, c.Offset)
	sort.Sort(system.SecondMenuSlice(secondMenuList))
	secondMenuDataJSON.SecondMenuList = secondMenuList
	secondMenuDataJSON.StartIndex = c.Offset

	if len(secondMenuList) == 0 {
		secondMenuDataJSON.Msg = "获取二级菜单失败"
	}

	c.GenerateJSON(secondMenuDataJSON)
}

func (c *GetController) GettelegramSecondMenus() {
	firstMenuUid := strings.TrimSpace(c.GetString("firMenuUid"))

	secondMenuDataJSON := new(datas.SecondMenuDataJSON)

	secondMenuList := system.GetSecondMenuListByFirstMenuUid(firstMenuUid)

	secondMenuDataJSON.Code = 200
	secondMenuDataJSON.SecondMenuList = secondMenuList
	c.GenerateJSON(secondMenuDataJSON)
}

func (c *GetController) GettelegramOneMenu() {
	menuUid := c.GetString("menuUid")

	dataJSON := new(datas.MenuDataJSON)
	menuInfo := system.GetMenuInfoByMenuUid(menuUid)
	if menuInfo.MenuUid == "" {
		dataJSON.Code = -1
		dataJSON.Msg = "该菜单项不存在"
	} else {
		dataJSON.MenuList = make([]system.MenuInfo, 0)
		dataJSON.MenuList = append(dataJSON.MenuList, menuInfo)
		dataJSON.Code = 200
	}
	c.Data["json"] = dataJSON
	_ = c.ServeJSONP()
}

func (c *GetController) GettelegramPowerItem() {
	powerID := c.GetString("powerID")
	powerItem := c.GetString("powerItem")

	params := make(map[string]string)
	params["power_uid__icontains"] = powerID
	params["power_item_icontains"] = powerItem

	l := system.GetPowerItemLenByMap(params)

	c.GetCutPage(l)

	powerItemDataJSON := new(datas.PowerItemDataJSON)
	powerItemDataJSON.DisplayCount = c.DisplayCount
	powerItemDataJSON.Code = 200
	powerItemDataJSON.CurrentPage = c.CurrentPage
	powerItemDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		powerItemDataJSON.PowerItemList = make([]system.PowerInfo, 0)
		c.GenerateJSON(powerItemDataJSON)
		return
	}

	powerItemDataJSON.StartIndex = c.Offset
	powerItemList := system.GetPowerItemByMap(params, c.DisplayCount, c.Offset)
	sort.Sort(system.PowerInfoSlice(powerItemList))
	powerItemDataJSON.PowerItemList = powerItemList

	c.GenerateJSON(powerItemDataJSON)
}

func (c *GetController) GettelegramRole() {
	//c.AllowCross()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", c.Ctx.Request.Header.Get("Origin"))
	//fmt.Print("2")
	roleName := strings.TrimSpace(c.GetString("roleName"))

	params := make(map[string]string)
	params["role_name__icontains"] = roleName

	l := system.GetRoleLenByMap(params)

	c.GetCutPage(l)

	roleInfoDataJSON := new(datas.RoleInfoDataJSON)
	roleInfoDataJSON.DisplayCount = c.DisplayCount
	roleInfoDataJSON.Code = 200
	roleInfoDataJSON.CurrentPage = c.CurrentPage
	roleInfoDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		roleInfoDataJSON.RoleInfoList = make([]system.RoleInfo, 0)
		c.GenerateJSON(roleInfoDataJSON)
		return
	}
	roleInfoDataJSON.StartIndex = c.Offset
	roleInfoDataJSON.RoleInfoList = system.GetRoleByMap(params, c.DisplayCount, c.Offset)

	c.GenerateJSON(roleInfoDataJSON)
}

func (c *GetController) GettelegramAllRole() {
	roleInfoDataJSON := new(datas.RoleInfoDataJSON)
	roleInfoList := system.GetRole()
	//fmt.Println(roleInfoList)
	if len(roleInfoList) == 0 {
		roleInfoDataJSON.Code = -1
	} else {
		roleInfoDataJSON.Code = 200
		roleInfoDataJSON.RoleInfoList = roleInfoList
	}
	c.GenerateJSON(roleInfoDataJSON)
}

func (c *GetController) GettelegramDeployTree() {
	roleUid := strings.TrimSpace(c.GetString("roleUid"))
	roleInfo := system.GetRoleByRoleUid(roleUid)

	allFirstMenu := system.GetMenuAll()
	sort.Sort(system.MenuInfoSlice(allFirstMenu))
	allSecondMenu := system.GetSecondMenuList()
	sort.Sort(system.SecondMenuSlice(allSecondMenu))
	allPower := system.GetPower()

	deployTreeJSON := new(datas.DeployTreeJSON)
	deployTreeJSON.Code = 200
	deployTreeJSON.AllFirstMenu = allFirstMenu
	deployTreeJSON.AllSecondMenu = allSecondMenu
	deployTreeJSON.AllPower = allPower
	deployTreeJSON.ShowFirstMenuUid = make(map[string]bool)
	for _, v := range strings.Split(roleInfo.ShowFirstUid, "||") {
		deployTreeJSON.ShowFirstMenuUid[v] = true
	}
	deployTreeJSON.ShowSecondMenuUid = make(map[string]bool)
	for _, v := range strings.Split(roleInfo.ShowSecondUid, "||") {
		deployTreeJSON.ShowSecondMenuUid[v] = true
	}
	deployTreeJSON.ShowPowerUid = make(map[string]bool)
	for _, v := range strings.Split(roleInfo.ShowPowerUid, "||") {
		deployTreeJSON.ShowPowerUid[v] = true
	}

	c.GenerateJSON(deployTreeJSON)
}

/*
* 获取操作员列表
 */
func (c *GetController) GettelegramOperator() {
	operatorName := strings.TrimSpace(c.GetString("operatorName"))

	params := make(map[string]string)
	params["user_id__icontains"] = operatorName

	l := user.GetOperatorLenByMap(params)
	c.GetCutPage(l)
	operatorDataJSON := new(datas.OperatorDataJSON)
	operatorDataJSON.DisplayCount = c.DisplayCount
	operatorDataJSON.Code = 200
	operatorDataJSON.CurrentPage = c.CurrentPage
	operatorDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		operatorDataJSON.OperatorList = make([]user.UserInfo, 0)
		c.GenerateJSON(operatorDataJSON)
		return
	}

	operatorDataJSON.StartIndex = c.Offset
	operatorDataJSON.OperatorList = user.GetOperatorByMap(params, c.DisplayCount, c.Offset)
	c.GenerateJSON(operatorDataJSON)
}

func (c *GetController) GettelegramOneOperator() {
	userId := strings.TrimSpace(c.GetString("userId"))

	userInfo := user.GetUserInfoByUserID(userId)

	operatorDataJSON := new(datas.OperatorDataJSON)
	operatorDataJSON.OperatorList = make([]user.UserInfo, 0)
	operatorDataJSON.OperatorList = append(operatorDataJSON.OperatorList, userInfo)

	operatorDataJSON.Code = 200

	c.GenerateJSON(operatorDataJSON)
}

func (c *GetController) GettelegramEditOperator() {
	userId := strings.TrimSpace(c.GetString("userId"))

	editOperatorDataJSON := new(datas.EditOperatorDataJSON)
	userInfo := user.GetUserInfoByUserID(userId)
	//fmt.Println(userInfo)
	editOperatorDataJSON.OperatorList = append(editOperatorDataJSON.OperatorList, userInfo)
	editOperatorDataJSON.RoleList = system.GetRole()
	editOperatorDataJSON.Code = 200

	c.GenerateJSON(editOperatorDataJSON)
}

func (c *GetController) GettelegramBankCard() {
	accountNameSearch := strings.TrimSpace(c.GetString("accountNameSearch"))
	params := make(map[string]string)
	params["account_name__icontains"] = accountNameSearch

	l := system.GetBankCardLenByMap(params)
	c.GetCutPage(l)

	bankCardDataJSON := new(datas.BankCardDataJSON)
	bankCardDataJSON.DisplayCount = c.DisplayCount
	bankCardDataJSON.Code = 200
	bankCardDataJSON.CurrentPage = c.CurrentPage
	bankCardDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		bankCardDataJSON.BankCardInfoList = make([]system.BankCardInfo, 0)
		c.GenerateJSON(bankCardDataJSON)
		return
	}

	bankCardDataJSON.StartIndex = c.Offset
	bankCardDataJSON.BankCardInfoList = system.GetBankCardByMap(params, c.DisplayCount, c.Offset)
	c.GenerateJSON(bankCardDataJSON)
}

func (c *GetController) GettelegramOneBankCard() {
	uid := strings.TrimSpace(c.GetString("uid"))
	bankCardInfo := system.GetBankCardByUid(uid)

	bankCardDataJSON := new(datas.BankCardDataJSON)
	bankCardDataJSON.Code = -1

	if bankCardInfo.Uid != "" {
		bankCardDataJSON.BankCardInfoList = append(bankCardDataJSON.BankCardInfoList, bankCardInfo)
		bankCardDataJSON.Code = 200
	}

	c.GenerateJSON(bankCardDataJSON)
}

/*
* 获取通道
 */

func (c *GetController) GettelegramRoad() {

	roadName := strings.TrimSpace(c.GetString("roadName"))
	productName := strings.TrimSpace(c.GetString("productName"))
	roadUid := strings.TrimSpace(c.GetString("roadUid"))
	payType := strings.TrimSpace(c.GetString("payType"))
	roadPoolCode := strings.TrimSpace(c.GetString("roadPoolCode"))
	paytypelei :=strings.TrimSpace(c.GetString("paytypelei"))

	params := make(map[string]string)
	params["road_name__icontains"] = roadName
	params["product_name__icontains"] = productName
	params["road_uid"] = roadUid
	params["pay_type"] = payType
	params["pay_typelei"] = paytypelei

	l := road.GetRoadLenByMap(params)
	c.GetCutPage(l)

	roadDataJSON := new(datas.RoadDataJSON)
	roadDataJSON.DisplayCount = c.DisplayCount
	roadDataJSON.Code = 200
	roadDataJSON.CurrentPage = c.CurrentPage
	roadDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		roadDataJSON.RoadInfoList = make([]road.RoadInfo, 0)
		c.GenerateJSON(roadDataJSON)
		return
	}


	roadDataJSON.StartIndex = c.Offset
	roadDataJSON.RoadInfoList,_= road.GetRoadInfoByMaptongji(params, c.DisplayCount, c.Offset,roadPoolCode)



	//roadDataJSON.RoadPool = road.GetRoadPoolByRoadPoolCode(roadPoolCode)
	c.GenerateJSON(roadDataJSON)
}

func (c *GetController) GettelegramAllRoad() {
	roadName := strings.TrimSpace(c.GetString("roadName"))
	params := make(map[string]string)
	params["road_name__icontains"] = roadName

	roadDataJSON := new(datas.RoadDataJSON)
	roadInfoList := road.GetAllRoad(params)

	roadDataJSON.Code = 200
	roadDataJSON.RoadInfoList = roadInfoList
	c.GenerateJSON(roadDataJSON)
}

/*
* 获取单个通道
 */
func (c *GetController) GettelegramOneRoad() {
	var u Urole
	json.Unmarshal(c.Ctx.Input.RequestBody, &u)
	//roadUid := strings.TrimSpace(u.RoadUid)
	roadUid := strings.TrimSpace(c.GetString("roadUid"))

	roadInfo := road.GetRoadInfoByRoadUid(roadUid)
	roadDataJSON := new(datas.RoadDataJSON)
	roadDataJSON.Code = -1

	if roadInfo.RoadUid != "" {
		roadDataJSON.RoadInfoList = append(roadDataJSON.RoadInfoList, roadInfo)
		roadDataJSON.Code = 200
	} else {
		roadDataJSON.RoadInfoList = make([]road.RoadInfo, 0)
	}

	c.GenerateJSON(roadDataJSON)
}

func (c *GetController) GettelegramRoadPool() {
	roadPoolName := strings.TrimSpace(c.GetString("roadPoolName"))
	roadPoolCode := strings.TrimSpace(c.GetString("roadPoolCode"))

	params := make(map[string]string)
	params["road_pool_name__icontains"] = roadPoolName
	params["road_pool_code__icontains"] = roadPoolCode

	l := road.GetRoadPoolLenByMap(params)
	c.GetCutPage(l)

	roadPoolDataJSON := new(datas.RoadPoolDataJSON)
	roadPoolDataJSON.DisplayCount = c.DisplayCount
	roadPoolDataJSON.Code = 200
	roadPoolDataJSON.CurrentPage = c.CurrentPage
	roadPoolDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		roadPoolDataJSON.RoadPoolInfoList = make([]road.RoadPoolInfo, 0)
		c.GenerateJSON(roadPoolDataJSON)
		return
	}


	roadPoolDataJSON.StartIndex = c.Offset
	roadPoolDataJSON.RoadPoolInfoList = road.GetRoadPoolByMap(params, c.DisplayCount, c.Offset)


	c.GenerateJSON(roadPoolDataJSON)
}

func (c *GetController) GettelegramAllRollPool() {
	rollPoolName := strings.TrimSpace(c.GetString("rollPoolName"))
	params := make(map[string]string)
	params["road_pool_name__icontains"] = rollPoolName

	roadPoolDataJSON := new(datas.RoadPoolDataJSON)
	roadPoolDataJSON.Code = 200
	roadPoolDataJSON.RoadPoolInfoList = road.GetAllRollPool(params)
	c.GenerateJSON(roadPoolDataJSON)
}

func (c *GetController) GettelegramMerchant() {
	merchantName := strings.TrimSpace(c.GetString("merchantName"))
	merchantNo := strings.TrimSpace(c.GetString("merchantNo"))

	params := make(map[string]string)
	params["merchant_name__icontains"] = merchantName
	params["merchant_uid__icontains"] = merchantNo

	l := merchant.GetMerchantLenByMap(params)
	c.GetCutPage(l)

	merchantDataJSON := new(datas.MerchantDataJSON)
	merchantDataJSON.DisplayCount = c.DisplayCount
	merchantDataJSON.Code = 200
	merchantDataJSON.CurrentPage = c.CurrentPage
	merchantDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		merchantDataJSON.MerchantList = make([]merchant.MerchantInfo, 0)
		c.GenerateJSON(merchantDataJSON)
		return
	}

	merchantDataJSON.StartIndex = c.Offset
	merchantDataJSON.MerchantList = merchant.GetMerchantListByMap(params, c.DisplayCount, c.Offset)
	merchantDataJSON.Accounts =Accounts.GetAccountByMap(params, c.DisplayCount, c.Offset)
	c.GenerateJSON(merchantDataJSON)
}

func (c *GetController) GettelegramAllMerchant() {
	merchantDataJSON := new(datas.MerchantDataJSON)
	merchantDataJSON.Code = 200
	merchantDataJSON.MerchantList = merchant.GetAllMerchant()
	c.GenerateJSON(merchantDataJSON)
}

func (c *GetController) GettelegramOneMerchant() {
	merchantUid := strings.TrimSpace(c.GetString("merchantUid"))
	merchantDataJSON := new(datas.MerchantDataJSON)

	if merchantUid == "" {
		merchantDataJSON.Code = -1
		c.GenerateJSON(merchantDataJSON)
		return
	}

	merchantInfo := merchant.GetMerchantByUid(merchantUid)

	merchantDataJSON.Code = 200
	merchantDataJSON.MerchantList = append(merchantDataJSON.MerchantList, merchantInfo)
	c.GenerateJSON(merchantDataJSON)
}

func (c *GetController) GettelegramOneMerchantDeploy() {
	merchantNo := strings.TrimSpace(c.GetString("merchantNo"))
	payType := strings.TrimSpace(c.GetString("payType"))

	merchantDeployDataJSON := new(datas.MerchantDeployDataJSON)

	merchantDeployInfo := merchant.GetMerchantDeployByUidAndPayType(merchantNo, payType)

	if merchantDeployInfo.Status == "active" {
		merchantDeployDataJSON.Code = 200
		merchantDeployDataJSON.MerchantDeploy = merchantDeployInfo
	} else {
		merchantDeployDataJSON.Code = -1
		merchantDeployDataJSON.MerchantDeploy = merchantDeployInfo
	}

	c.GenerateJSON(merchantDeployDataJSON)
}

func (c *GetController) GettelegramAllAccount() {
	accountDataJSON := new(datas.AccountDataJSON)
	accountDataJSON.Code = 200

	accountDataJSON.AccountList = Accounts.GetAllAccount()

	c.GenerateJSON(accountDataJSON)
}

func (c *GetController) GettelegramAccount() {
	accountName := strings.TrimSpace(c.GetString("accountName"))
	accountUid := strings.TrimSpace(c.GetString("accountNo"))

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
	accountDataJSON.AccountList = Accounts.GetAccountByMap(params, c.DisplayCount, c.Offset)
	c.GenerateJSON(accountDataJSON)
}

func (c *GetController) GettelegramAccounttody() {
	accountName := strings.TrimSpace(c.GetString("accountName"))
	accountUid := strings.TrimSpace(c.GetString("accountNo"))
	start := strings.TrimSpace(c.GetString("startTime"))
	end := strings.TrimSpace(c.GetString("endTime"))
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
	accountDataJSON.AccountTody = Accounts.GetAccountByMaptody(params, c.DisplayCount, c.Offset,start,end,accountName)
	c.GenerateJSON(accountDataJSON)
}



func (c *GetController) GettelegramAccounttodytongdao() {
	//accountName := strings.TrimSpace(c.GetString("accountName"))
	accountUid := strings.TrimSpace(c.GetString("formuid"))
	//start := strings.TrimSpace(c.GetString("startTime"))
	//end := strings.TrimSpace(c.GetString("endTime"))
	params := make(map[string]string)
	//params["account_name__icontains"] = accountName
	params["merchant_uid"] = accountUid



	l := merchant.GetMerchanttDeployLenByMap(params)
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
	accountDataJSON.AccountTody = merchant.GetMerchantListDeployByMap(params, c.DisplayCount, c.Offset)
	c.GenerateJSON(accountDataJSON)
}

func (c *GetController) GettelegramOneAccount() {
	//从http的body中获取accountUid字段，并且这个字段是string类型
	accountUid := strings.TrimSpace(c.GetString("accountUid"))
	//new一个accountDataJSON结构体对象，用来做jsonp返回
	accountDataJSON := new(datas.AccountDataJSON)
	//用accountuid作为过滤字段，从数据库中读取一条信息
	accountInfo := Accounts.GetAccountByUid(accountUid)
	//code初始值为200
	accountDataJSON.Code = 200
	//将从数据库读出来的数据插入到accountList数组中
	accountDataJSON.AccountList = append(accountDataJSON.AccountList, accountInfo)
	//返回jsonp格式的数据给前端
	c.GenerateJSON(accountDataJSON)
}

func (c *GetController) GettelegramAccountHistory() {
	accountName := strings.TrimSpace(c.GetString("accountHistoryName"))
	accountUid := strings.TrimSpace(c.GetString("accountHistoryNo"))
	operatorType := strings.TrimSpace(c.GetString("operatorType"))
	operatorTypep := strings.TrimSpace(c.GetString("operatorTypep"))
	startTime := c.GetString("startTime")
	endTime := c.GetString("endTime")
	switch operatorType {
	case "plus-amount":
		operatorType = common.PLUS_AMOUNT
	case "sub-amount":
		operatorType = common.SUB_AMOUNT
	case "freeze-amount":
		operatorType = common.FREEZE_AMOUNT
	case "unfreeze-amount":
		operatorType = common.UNFREEZE_AMOUNT
	}

	if (operatorTypep=="tiaozhen"){
		operatorTypep="1"
	}

	params := make(map[string]string)
	params["account_name__icontains"] = accountName
	params["account_uid__icontains"] = accountUid
	params["type"] = operatorType
	params["notify"] = operatorTypep
	params["create_time__gte"] = startTime
	params["create_time__lte"] = endTime

	l := Accounts.GetAccountHistoryLenByMap(params)
	c.GetCutPage(l)

	accountHistoryDataJSON := new(datas.AccountHistoryDataJSON)
	accountHistoryDataJSON.DisplayCount = c.DisplayCount
	accountHistoryDataJSON.Code = 200
	accountHistoryDataJSON.CurrentPage = c.CurrentPage
	accountHistoryDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		accountHistoryDataJSON.AccountHistoryList = make([]Accounts.AccountHistoryInfo, 0)
		c.GenerateJSON(accountHistoryDataJSON)
		return
	}

	accountHistoryDataJSON.StartIndex = c.Offset
	accountHistoryDataJSON.AccountHistoryList = Accounts.GetAccountHistoryByMap(params, c.DisplayCount, c.Offset)
	c.GenerateJSON(accountHistoryDataJSON)
}


func (c *GetController) GettelegramAccountyunyingHistory() {
	accountName := strings.TrimSpace(c.GetString("accountHistoryName"))
	accountUid := strings.TrimSpace(c.GetString("accountHistoryNo"))
	operatorType := strings.TrimSpace(c.GetString("operatorType"))
	startTime := c.GetString("startTime")
	endTime := c.GetString("endTime")

	switch operatorType {
	case "plus-amount":
		operatorType = common.PLUS_AMOUNT
	case "sub-amount":
		operatorType = common.SUB_AMOUNT
	case "freeze-amount":
		operatorType = common.FREEZE_AMOUNT
	case "unfreeze-amount":
		operatorType = common.UNFREEZE_AMOUNT
	}
	params := make(map[string]string)
	params["account_name__icontains"] = accountName
	params["account_uid__icontains"] = accountUid
	params["type"] = operatorType
	params["create_time__gte"] = startTime
	params["create_time__lte"] = endTime

	l := mauser.GetAccountHistoryLenByMap(params)
	c.GetCutPage(l)

	accountHistoryDataJSON := new(datas.MaunsterlistHistoryDataJSON)
	accountHistoryDataJSON.DisplayCount = c.DisplayCount
	accountHistoryDataJSON.Code = 200
	accountHistoryDataJSON.CurrentPage = c.CurrentPage
	accountHistoryDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		accountHistoryDataJSON.AccountHistoryList = make([]mauser.MauserlistHistoryInfo, 0)
		c.GenerateJSON(accountHistoryDataJSON)
		return
	}

	accountHistoryDataJSON.StartIndex = c.Offset
	accountHistoryDataJSON.AccountHistoryList = mauser.GetAccountHistoryByMap(params, c.DisplayCount, c.Offset)
	c.GenerateJSON(accountHistoryDataJSON)
}

func (c *GetController) Delecttabletelegram() {

	payWayCode := strings.TrimSpace(c.GetString("payWayCode"))
	startTime := c.GetString("startTime")
	endTime := c.GetString("endTime")

	//fmt.Print(payWayCode)

	params := make(map[string]string)
	params["table"] = payWayCode
	params["start"] = startTime
	params["endstart"] = endTime

	l := mauser.Deleteorm(params)
	//fmt.Print(l)

	accountHistoryDataJSON := new(datas.MaunsterlistHistoryDataJSON)
	accountHistoryDataJSON.DisplayCount = c.DisplayCount
	accountHistoryDataJSON.Code = 200
	accountHistoryDataJSON.Msg = l

	c.GenerateJSON(accountHistoryDataJSON)
}

func (c *GetController) DhosttelegramMerchantdefault() {

	dispayid := c.GetString("dispayid")
	status := c.GetString("status")



	_=merchant.UpdateStautskphone(status, dispayid)
	//fmt.Print(l)

	accountHistoryDataJSON := new(datas.MaunsterlistHistoryDataJSON)
	accountHistoryDataJSON.DisplayCount = c.DisplayCount
	accountHistoryDataJSON.Code = 200
	accountHistoryDataJSON.Msg = "更新成功"

	c.GenerateJSON(accountHistoryDataJSON)
}

func (c *GetController) GettelegramAgent() {
	agentName := strings.TrimSpace(c.GetString("agentName"))
	params := make(map[string]string)
	params["agnet_name__icontains"] = agentName

	l := agent.GetAgentInfoLenByMap(params)
	c.GetCutPage(l)

	agentDataJSON := new(datas.AgentDataJSON)
	agentDataJSON.DisplayCount = c.DisplayCount
	agentDataJSON.Code = 200
	agentDataJSON.CurrentPage = c.CurrentPage
	agentDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		agentDataJSON.AgentList = make([]agent.AgentInfo, 0)
		c.GenerateJSON(agentDataJSON)
		return
	}

	agentDataJSON.StartIndex = c.Offset
	agentDataJSON.AgentList = agent.GetAgentInfoByMap(params, c.DisplayCount, c.Offset)
	c.GenerateJSON(agentDataJSON)
}

func (c *GetController) GettelegramAllAgent() {
	agentName := strings.TrimSpace(c.GetString("agentName"))
	params := make(map[string]string)
	params["agent_name__icontains"] = agentName

	agentDataJSON := new(datas.AgentDataJSON)
	agentDataJSON.Code = 200
	agentDataJSON.AgentList = agent.GetAllAgentByMap(params)

	c.GenerateJSON(agentDataJSON)
}

func (c *GetController) GettelegramProduct() {
	supplierCode2Name := common.GetSupplierMap()
	productDataJSON := new(datas.ProductDataJSON)
	productDataJSON.Code = 200
	productDataJSON.ProductMap = supplierCode2Name
	c.GenerateJSON(productDataJSON)
}

func (c *GetController) GettelegramProductshanghu() {
	supplierCode2Name := merchant.GetAllMerchant()
	productDataJSON := new(datas.ProductDataJSON)
	productDataJSON.Code = 200
	productDataJSON.Merchant = supplierCode2Name
	c.GenerateJSON(productDataJSON)
}

func (c *GetController) GettelegramProductlun() {
	roadPoolDataJSON := new(datas.RoadPoolDataJSON)
	roadPoolDataJSON.RoadPoolInfoList = road.GetRoadPoollist()
	roadPoolDataJSON.Code = 200

	c.GenerateJSON(roadPoolDataJSON)
}

func (c *GetController) GettelegramAgentToMerchant() {
	agentUid := strings.TrimSpace(c.GetString("agentUid"))
	merchantUid := strings.TrimSpace(c.GetString("merchantUid"))

	params := make(map[string]string)
	params["belong_agent_uid"] = agentUid
	params["merchant_uid"] = merchantUid

	l := merchant.GetMerchantLenByParams(params)
	c.GetCutPage(l)

	merchantDataJSON := new(datas.MerchantDataJSON)
	merchantDataJSON.DisplayCount = c.DisplayCount
	merchantDataJSON.Code = 200
	merchantDataJSON.CurrentPage = c.CurrentPage
	merchantDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		merchantDataJSON.MerchantList = make([]merchant.MerchantInfo, 0)
	} else {
		merchantDataJSON.MerchantList = merchant.GetMerchantByParams(params, c.DisplayCount, c.Offset)
	}

	c.GenerateJSON(merchantDataJSON)
}

/*
* 获取订单数据
 */
func (c *GetController) GettelegramOrder() {
	startTime := strings.TrimSpace(c.GetString("startTime"))
	endTime := strings.TrimSpace(c.GetString("endTime"))
	merchantName := strings.TrimSpace(c.GetString("merchantName"))
	orderNo := strings.TrimSpace(c.GetString("merchantOrderId"))
	bankOrderId := strings.TrimSpace(c.GetString("bankOrderId"))
	status := strings.TrimSpace(c.GetString("orderStatus"))
	supplierUid := strings.TrimSpace(c.GetString("supplierUid"))
	payWayCode := strings.TrimSpace(c.GetString("payWayCode"))
	freeStatus := strings.TrimSpace(c.GetString("freeStatus"))
	merchantid := strings.TrimSpace(c.GetString("merchantid"))
	mauser := strings.TrimSpace(c.GetString("mauser"))
	userPhone := strings.TrimSpace(c.GetString("userPhone"))

	params := make(map[string]string)
	params["create_time__gte"] = startTime
	params["create_time__lte"] = endTime
	params["merchant_name__icontains"] = merchantName
	params["merchant_order_id"] = orderNo
	params["bank_order_id"] = bankOrderId
	params["merchant_uid"] = merchantid
	params["status"] = status
	params["roll_pool_code"] = supplierUid
	params["road_uid"] = payWayCode
	params["mauser"] = mauser
	params["user_phone"] = userPhone

	switch freeStatus {
	case "free":
		params["free"] = "yes"
	case "unfree":
		params["unfree"] = "yes"
	case "refund":
		params["refund"] = "yes"
	}



	l := order.GetOrderLenByMap(params)
	c.GetCutPage(l)

	orderDataJSON := new(datas.OrderDataJSON)
	orderDataJSON.DisplayCount = c.DisplayCount
	orderDataJSON.Code = 200
	orderDataJSON.CurrentPage = c.CurrentPage
	orderDataJSON.TotalPage = c.TotalPage
	orderDataJSON.Count = l

	orderDataJSON.AllAmount,orderDataJSON.Succnum,orderDataJSON.SuccAmount,orderDataJSON.Sxfpress,orderDataJSON.Sxfallmount = order.GetAllAmountByMap(params)
	if c.Offset < 0 {
		orderDataJSON.OrderList = make([]order.OrderInfo, 0)
		c.GenerateJSON(orderDataJSON)
		return
	}

	orderDataJSON.StartIndex = c.Offset
	orderDataJSON.OrderListv = order.GetOrderByMap(params, c.DisplayCount, c.Offset)
	orderDataJSON.SuccessRate = order.GetSuccessRateByMap(params)
	params["status"] = common.SUCCESS

	c.GenerateJSON(orderDataJSON)
}

func (c *GetController) GettelegramOneOrder() {
	bankOrderId := strings.TrimSpace(c.GetString("bankOrderId"))
	orderDataJSON := new(datas.OrderDataJSON)
	orderInfo := order.GetOneOrder(bankOrderId)

	orderDataJSON.Code = 200
	orderDataJSON.OrderList = append(orderDataJSON.OrderList, orderInfo)
	notifyInfo := notify.GetNotifyInfoByBankOrderId(bankOrderId)
	if notifyInfo.Url == "" || len(notifyInfo.Url) == 0 {
		orderDataJSON.NotifyUrl = orderInfo.NotifyUrl
	} else {
		orderDataJSON.NotifyUrl = notifyInfo.Url
	}
	c.GenerateJSON(orderDataJSON)
}

func (c *GetController) GettelegramOrderProfit() {
	startTime := strings.TrimSpace(c.GetString("startTime"))
	endTime := strings.TrimSpace(c.GetString("endTime"))
	merchantName := strings.TrimSpace(c.GetString("merchantName"))
	agentName := strings.TrimSpace(c.GetString("agentName"))
	bankOrderId := strings.TrimSpace(c.GetString("bankOrderId"))
	status := strings.TrimSpace(c.GetString("orderStatus"))
	supplierUid := strings.TrimSpace(c.GetString("supplierUid"))
	payWayCode := strings.TrimSpace(c.GetString("payWayCode"))
	merchantid := strings.TrimSpace(c.GetString("merchantid"))

	params := make(map[string]string)
	params["create_time__gte"] = startTime
	params["create_time__lte"] = endTime
	params["merchant_name__icontains"] = merchantName
	params["agent_name__icontains"] = agentName
	params["bank_order_id"] = bankOrderId
	params["merchant_uid"] = merchantid
	params["status"] = status
	params["roll_pool_code"] = supplierUid
	params["pay_type_code"] = payWayCode

	l := order.GetOrderProfitLenByMap(params)
	c.GetCutPage(l)

	listDataJSON := new(datas.ListDataJSON)
	listDataJSON.DisplayCount = c.DisplayCount
	listDataJSON.Code = 200
	listDataJSON.CurrentPage = c.CurrentPage
	listDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		listDataJSON.List = make([]order.OrderProfitInfo, 0)
		c.GenerateJSON(listDataJSON)
		return
	}
	AllAmount:=0.00
	UserInAmount:=0.00
	listDataJSON.StartIndex = c.Offset
	var  vcount interface{}
	var  success interface{}
	AllAmount,UserInAmount,listDataJSON.PlatformProfit,listDataJSON.AgentProfit,listDataJSON.SupplierProfit= order.GetAllAmountByMapinfo(params)

	//fmt.Print(AllAmount)
	listDataJSON.AllAmount= AllAmount
	listDataJSON.UserInAmount =UserInAmount
	listDataJSON.List, vcount,success = order.GetOrderProfitByMap(params, c.DisplayCount, c.Offset)
	supplierAll := 0.0
	platformAll := 0.0
	agentAll := 0.0
	allAmount := 0.0
	userInamount:=0.0

	for _, v := range listDataJSON.Listall {

		allAmount += v.FactAmount
		supplierAll += v.SupplierProfit
		platformAll += v.PlatformProfit
		agentAll += v.AgentProfit
		userInamount += v.UserInAmount
	}



	listDataJSON.Count = vcount
	listDataJSON.Successface =success
	//listDataJSON.UserInAmount, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", userInamount), 3)
	//listDataJSON.SupplierProfit, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", supplierAll), 3)
	//listDataJSON.PlatformProfit, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", platformAll), 3)
	listDataJSON.AgentProfit, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", agentAll), 3)
	//listDataJSON.AllAmount, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", allAmount), 3)

	c.GenerateJSON(listDataJSON)
}



func (c *GetController) GettelegramOrderProfitcha() {
	startTime := strings.TrimSpace(c.GetString("startTime"))
	endTime := strings.TrimSpace(c.GetString("endTime"))
	merchantName := strings.TrimSpace(c.GetString("merchantName"))
	agentName := strings.TrimSpace(c.GetString("agentName"))
	bankOrderId := strings.TrimSpace(c.GetString("bankOrderId"))
	status := strings.TrimSpace(c.GetString("orderStatus"))
	supplierUid := strings.TrimSpace(c.GetString("supplierUid"))
	payWayCode := strings.TrimSpace(c.GetString("payWayCode"))

	params := make(map[string]string)
	params["create_time__gte"] = startTime
	params["create_time__lte"] = endTime
	params["merchant_name__icontains"] = merchantName
	params["agent_name__icontains"] = agentName
	params["bank_order_id"] = bankOrderId
	params["status"] = status
	params["roll_pool_code"] = supplierUid
	params["pay_type_code"] = payWayCode



	listDataJSON := new(datas.ListDataJSON)
	listDataJSON.DisplayCount = c.DisplayCount
	listDataJSON.Code = 200
	listDataJSON.CurrentPage = c.CurrentPage
	listDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		listDataJSON.List = make([]order.OrderProfitInfo, 0)
		c.GenerateJSON(listDataJSON)
		return
	}

	listDataJSON.StartIndex = c.Offset
	var  vcount interface{}
	var  success interface{}

	listDataJSON.List,vcount,success = order.GetOrderProfitByMap(params, c.DisplayCount, c.Offset)
	supplierAll := 0.0
	platformAll := 0.0
	agentAll := 0.0
	allAmount := 0.0
	userInamount:=0.0

	for _, v := range listDataJSON.Listall {

		allAmount += v.FactAmount
		supplierAll += v.SupplierProfit
		platformAll += v.PlatformProfit
		agentAll += v.AgentProfit
		userInamount += v.UserInAmount
	}

	listDataJSON.Count = vcount
	listDataJSON.Successface =success
	listDataJSON.UserInAmount, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", userInamount), 3)
	listDataJSON.SupplierProfit, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", supplierAll), 3)
	listDataJSON.PlatformProfit, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", platformAll), 3)
	listDataJSON.AgentProfit, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", agentAll), 3)
	listDataJSON.AllAmount, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", allAmount), 3)

	c.GenerateJSON(listDataJSON)
}




func (c *GetController) GettelegramPayFor() {
	startTime := strings.TrimSpace(c.GetString("startTime"))
	endTime := strings.TrimSpace(c.GetString("endTime"))
	merchantOrderId := strings.TrimSpace(c.GetString("accountHistoryName"))
	bankOrderId := strings.TrimSpace(c.GetString("accountHistoryNo"))
	marchantorderid := strings.TrimSpace(c.GetString("marchantorderid"))
	bankorderid := strings.TrimSpace(c.GetString("bankorderid"))


	status := strings.TrimSpace(c.GetString("operatorType"))

	params := make(map[string]string)
	params["create_time__lte"] = endTime
	params["create_time__gte"] = startTime
	params["merchant_name"] = merchantOrderId
	params["bank_order_id"] = marchantorderid
	params["merchant_order_id"] = bankorderid
	params["merchant_uid"] = bankOrderId

	params["status"] = status

	l := payfor.GetPayForLenByMap(params)
	c.GetCutPage(l)

	listDataJSON := new(datas.PayForDataJSON)
	listDataJSON.DisplayCount = c.DisplayCount
	listDataJSON.Code = 200
	listDataJSON.CurrentPage = c.CurrentPage
	listDataJSON.TotalPage = c.TotalPage

	if c.Offset < 0 {
		listDataJSON.PayForList = make([]payfor.PayforInfo, 0)
		c.GenerateJSON(listDataJSON)
		return
	}

	listDataJSON.StartIndex = c.Offset
	listDataJSON.PayForList = payfor.GetPayForByMap(params, c.DisplayCount, c.Offset)
	listDataJSON.TotalAmount,listDataJSON.TotalInt,listDataJSON.TotalsucAmount,listDataJSON.Payforfee=payfor.GetAllAmountByMap(params)


	for index, p := range listDataJSON.PayForList {
		if p.MerchantName == "" {
			listDataJSON.PayForList[index].MerchantName = "任意下发"
		}
		if p.MerchantOrderId == "" {
			listDataJSON.PayForList[index].MerchantOrderId = "任意发下"
		}
		if p.RoadName == "" {
			listDataJSON.PayForList[index].RoadName = "无"
		}
	}
	c.GenerateJSON(listDataJSON)
}

func (c *GetController) GettelegramOnePayFor() {
	bankOrderId := strings.TrimSpace(c.GetString("bankOrderId"))

	payForInfo := payfor.GetPayForByBankOrderId(bankOrderId)

	listDataJSON := new(datas.PayForDataJSON)
	listDataJSON.Code = 200
	listDataJSON.PayForList = append(listDataJSON.PayForList, payForInfo)

	c.GenerateJSON(listDataJSON)
}

func (c *GetController) GettelegramBalance() {
	/*roadName := strings.TrimSpace(c.GetString("roadName"))
	roadUid := strings.TrimSpace(c.GetString("roadUid"))*/

	/*var roadInfo road.RoadInfo
	if roadUid != "" {
		roadInfo = road.GetRoadInfoByRoadUid(roadUid)
	} else {
		roadInfo = road.GetRoadInfoByName(roadName)
	}*/

	balanceDataJSON := new(datas.BalanceDataJSON)
	balanceDataJSON.Code = 200

	/*supplier := controller.GetPaySupplierByCode(roadInfo.ProductUid)
	if supplier == nil {
		balanceDataJSON.Code = -1
		balanceDataJSON.Balance = -1.00
	} else {
		balance := supplier.BalanceQuery(roadInfo)
		balanceDataJSON.Balance = balance
	}*/
	// TODO 从gateway获取账户余额
	balanceDataJSON.Balance = 1

	c.GenerateJSON(balanceDataJSON)
}

func (c *GetController) GettelegramNotifyBankOrderIdList() {
	startTime := strings.TrimSpace(c.GetString("startTime"))
	endTime := strings.TrimSpace(c.GetString("endTime"))
	merchantUid := strings.TrimSpace(c.GetString("merchantUid"))
	notifyType := strings.TrimSpace(c.GetString("notifyType"))

	params := make(map[string]string)
	params["create_time__gte"] = startTime
	params["create_time_lte"] = endTime
	params["merchant_uid"] = merchantUid
	params["type"] = notifyType

	bankOrderIdListJSON := new(datas.NotifyBankOrderIdListJSON)
	bankOrderIdListJSON.Code = 200
	bankOrderIdListJSON.NotifyIdList = notify.GetNotifyBankOrderIdListByParams(params)
	c.GenerateJSON(bankOrderIdListJSON)
}

/*
* 获取利润表
 */
func (c *GetController) GettelegramProfit() {
	merchantUid := strings.TrimSpace(c.GetString("merchantUid"))
	agentUid := strings.TrimSpace(c.GetString("agentUid"))
	supplierUid := strings.TrimSpace(c.GetString("supplierUid"))
	payType := strings.TrimSpace(c.GetString("payType"))
	startTime := strings.TrimSpace(c.GetString("startTime"))
	endTime := strings.TrimSpace(c.GetString("endTime"))

	params := make(map[string]string)
	params["merchant_uid"] = merchantUid
	params["agent_uid"] = agentUid
	params["roll_pool_code"] = supplierUid
	params["pay_type_code"] = payType
	params["create_time__gte"] = startTime
	params["create_time__lte"] = endTime

	profitListJSON := new(datas.ProfitListJSON)
	profitListJSON.Code = 200




	profitListJSON.ProfitList = order.GetPlatformProfitByMap(params)
	TotalAmountyear := order.GetPlatformProfitByMapfy()
	profitListJSON.TotalAmount = 0.00
	profitListJSON.PlatformTotalProfit = 0.00
	profitListJSON.AgentTotalProfit = 0.00


	for i, p := range profitListJSON.ProfitList {
		profitListJSON.TotalAmount += p.OrderAmount
		profitListJSON.PlatformTotalProfit += p.PlatformProfit
		profitListJSON.AgentTotalProfit += p.AgentProfit
		profitListJSON.TodaySubmit, profitListJSON.TodaySuccess, profitListJSON.TodayProfit,profitListJSON.TodayTotal, profitListJSON.YesterdaySubmit, profitListJSON.YesterdaySuccess, profitListJSON.YesterdayProfit, profitListJSON.YesterdayTotal =order.GetPlatformToday(p.MerchantUid,p.PayTypeCode)

		if p.AgentName == "" {
			p.AgentName = "无代理商"
		}
		profitListJSON.ProfitList[i] = p
	}
	profitListJSON.AgentTotalProfitutf = TotalAmountyear
	c.GenerateJSON(profitListJSON)
}

func (c *GetController) Gettelegramzhuzhuang() {
	var datalist []interface{}
	t := time.Now()

	profitListJSON := new(datas.Pzhuzhuang)
	profitListJSON.Code = 200


	v:=getYearMonthToDay(t.Year(), int(t.Month()))
	for i:=0;i<v;i++{
		firstDayx := time.Date(t.Year(), t.Month(), 1+i, 0, 0, 0, 0, time.UTC).Format("2006-01-02 00:00:00")
		lastDay :=  time.Date(t.Year(), t.Month(), 1+i, 23, 59, 59, 0, time.Local).Format("2006-01-02 15:04:05")
		profitListJSON.TotalAmountyear=order.GetPlatformProfitzhuzhuang(firstDayx,lastDay)
		for _,p:=range profitListJSON.TotalAmountyear{
			//profitListJSON.TotalAmountyear[i] = p.OrderAmount
			profitzz := order.PlatformProfitzz{Data: firstDayx,OrderAmount:p.OrderAmount,OrderCount:p.OrderCount,PlatformProfit:p.PlatformProfit,AgentProfit:p.AgentProfit}
			datalist=append(datalist, profitzz)
		}
		//profitzz := order.PlatformProfitzz{Data: firstDayx,OrderAmount:p.A}


	}
	profitListJSON.Code=200
	profitListJSON.Msg="获取成功"
	profitListJSON.Totalv = datalist
	c.GenerateJSON(profitListJSON)
}
func gettelegramYearMonthToDay(year int, month int) int {
	// 有31天的月份
	day31 := map[int]struct{}{
		1:  struct{}{},
		3:  struct{}{},
		5:  struct{}{},
		7:  struct{}{},
		8:  struct{}{},
		10: struct{}{},
		12: struct{}{},
	}
	if _, ok := day31[month]; ok {
		return 31
	}
	// 有30天的月份
	day30 := map[int]struct{}{
		4:  struct{}{},
		6:  struct{}{},
		9:  struct{}{},
		11: struct{}{},
	}
	if _, ok := day30[month]; ok {
		return 30
	}
	// 计算是平年还是闰年
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		// 得出2月的天数
		return 29
	}
	// 得出2月的天数
	return 28
}


