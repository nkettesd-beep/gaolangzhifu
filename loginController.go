package controllers

import (
	"boss/common"
	"boss/datas"
	"boss/models/user"
	"boss/utils"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/adapter/validation"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"strconv"
	"strings"

	"time"
)

type LoginController struct {
	beego.Controller
}
type U struct {
	UserID string
	Passwd string
	jwt.StandardClaims
	Code   string

}
type Ut struct {
	Token string
}
func (c *LoginController) Prepare() {

}

//跨域
//func (c *LoginController) AllowCross() {
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")       //允许访问源
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")    //允许post访问
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
//	c.Ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
//}



func (c *LoginController) Login() {
	//c.AllowCross()
	//c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", c.Ctx.Request.Header.Get("Origin"))

	var u U
	json.Unmarshal(c.Ctx.Input.RequestBody, &u)
	userID := u.UserID
	passWD :=u.Passwd
	code := u.Code
	//userID := c.GetString("userID")
	//passWD :=c.GetString("passwd")
	//code := c.GetString("Code")

	dataJSON := new(datas.KeyDataJSON)

	valid := validation.Validation{}

	if v := valid.Required(userID, "userID"); !v.Ok {
		dataJSON.Key = v.Error.Key
		dataJSON.Key = "手机号不能为空！"
	} else if v := valid.Required(passWD, "passWD"); !v.Ok {
		dataJSON.Key = v.Error.Key
		dataJSON.Msg = "登录密码不能为空！"
	}
	//else if v := valid.Length(code, common.VERIFY_CODE_LEN, "code"); !v.Ok {
	//	dataJSON.Key = v.Error.Key
	//	dataJSON.Msg = "验证码不正确！"
	//}

	secret:=user.GetUserInfoByUserID(userID)
	if (secret.Googlekey !=""){
		b,_:= strconv.ParseInt(code,10,32)
		bc := VerifyCode(secret.Googlekey, int32(b))

		if bc {
			//fmt.Println("验证成功！")
			//dataJSON.Key = "code"
			//dataJSON.Msg = "验证成功！"
		} else {
			//fmt.Println("验证失败！")
			dataJSON.Key = common.UNACTIVE
			dataJSON.Msg = "GOOGLE验证码不正确！"
		}
	}




	userInfo := user.GetUserInfoByUserID(userID)

	if userInfo.UserId == "" {
		dataJSON.Key = "userID"
		dataJSON.Msg = "用户不存在，请求联系管理员！"
	} else {
		//codeInterface := c.GetSession("verifyCode")
		if userInfo.Passwd != utils.GetMD5Upper(passWD) {
			dataJSON.Key = "passWD"
			dataJSON.Msg = "密码不正确！"
		//} else if codeInterface == nil {
		//	dataJSON.Key = "code"
		//	dataJSON.Msg = "验证码失效！"
		//} else if code != codeInterface.(string) {
		//	dataJSON.Key = "code"
		//	dataJSON.Msg = "验证码不正确！"
		} else if userInfo.Status == common.UNACTIVE {
			dataJSON.Key = common.UNACTIVE
			dataJSON.Msg = "用户已被冻结！"
		} else if userInfo.Status == "del" {
			dataJSON.Key = "del"
			dataJSON.Msg = "用户已被删除！"
		}
	}

	go func() {
		userInfo.Ip = c.Ctx.Input.IP()
		user.UpdateUserInfoIP(userInfo)
	}()



	if dataJSON.Key == "" {
		_ = c.SetSession("userID", userID)
		_= c.DelSession("verifyCode")
		dataJSON.Code=200
		token,_:=CreateToken(userID)
		dataJSON.Data.Token=token
	}
	//dataJSON.Code=200
	//dataJSON.Data.Token="admin"

	c.Data["json"] = dataJSON
	_ = c.ServeJSON()
}

func  (c *LoginController) Googleyanzheng() {

	secret := GetSecret()
	fmt.Println("secret:" + secret)
	dataJSON := new(datas.BaseDataJSON)
	dataJSON.Code=200
	dataJSON.Yzheng=secret
	c.Data["json"] = dataJSON
	_ = c.ServeJSON()


}



func GetSecret() string {
	randomStr := randStr(16)
	return strings.ToUpper(randomStr)
}

func randStr(strSize int) string {
	dictionary := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, strSize)
	_, _ = rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

// 为了考虑时间误差，判断前当前时间及前后30秒时间
func VerifyCode(secret string, code int32) bool {
	// 当前google值
	if getCode(secret, 0) == code {
		return true
	}

	// 前30秒google值
	if getCode(secret, -30) == code {
		return true
	}

	// 后30秒google值
	if getCode(secret, 30) == code {
		return true
	}

	return false
}

// 获取Google Code
func getCode(secret string, offset int64) int32 {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix() + offset
	return int32(oneTimePassword(key, toBytes(epochSeconds/30)))
}

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}






/*
* 退出登录,删除session中的数据，避免数据量过大，内存吃紧
 */

func (c *LoginController) Logout() {
	dataJSON := new(datas.BaseDataJSON)

	_ = c.DelSession("userID")
	dataJSON.Code = 200
	c.Data["json"] = dataJSON
	_ = c.ServeJSON()

}

const secretkey = "131486555aDFFDSdsbcdefgajf"
func CreateToken(user string) (string, error) {
	customClaims := &U{
		UserID: user,
		StandardClaims: jwt.StandardClaims{
			Audience:  "支付项目",                                             //颁发给谁，就是使用的一方
			ExpiresAt: time.Now().Add(time.Duration(time.Hour * 1)).Unix(), //过期时间
			//Id:        "",//非必填
			IssuedAt: time.Now().Unix(), //颁发时间
			Issuer:   "支付系统",             //颁发者
			//NotBefore: 0,         //生效时间，就是在这个时间前不能使用，非必填
			Subject: "17665664",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(secretkey))
	if err != nil {
		return "", err
	}
	//defer func() {
	//    e := recover()
	//    if e != nil {
	//        panic(e)
	//    }
	//}()

	return tokenString, nil
}

/*
* 验证码获取，如果获取成功，并将验证码存到session中
 */
func (c *LoginController) GetVerifyImg() {
	Image, verifyCode := utils.GenerateVerifyCodeImg()
	if Image == nil || len(verifyCode) != common.VERIFY_CODE_LEN {
		logs.Error("获取验证码图片失败！")
	} else {
		_ = c.SetSession("verifyCode", verifyCode)
	}
	logs.Info("验证码：", verifyCode)
	if Image == nil {
		logs.Error("生成验证码失败！")
	} else {
		_, _ = Image.WriteTo(c.Ctx.ResponseWriter)
	}
}


func (c *LoginController) Userinfo() {
	dataJSON := new(datas.KeyDataJSON)
	var u Ut
	json.Unmarshal(c.Ctx.Input.RequestBody, &u)
	token := u.Token
	_, err := ParseToken(token)
	if err == nil {
		dataJSON.Code=200
		dataJSON.Data.Token=token
		c.Data["json"] = dataJSON
		_ = c.ServeJSON()

	}else {
		dataJSON.Code=404
		dataJSON.Data.Token=token
		c.Data["json"] = dataJSON
		_ = c.ServeJSON()
	}


}
func ParseToken(tokenString string) (*U, error) {
	/*
		jwt.ParseWithClaims有三个参数
		第一个就是加密后的token字符串
		第二个是Claims, &CustomClaims 传指针的目的是。然让它在我这个里面写数据
		第三个是一个自带的回调函数，将秘钥和错误return出来即可，要通过密钥解密
	*/
	token, err := jwt.ParseWithClaims(tokenString, &U{}, func(token *jwt.Token) (interface{}, error) {
		//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//    return nil, fmt.Errorf("算法类型: %v", token.Header["alg"])
		//}
		return []byte(secretkey), nil
	})
	if claims, ok := token.Claims.(*U); ok && token.Valid {
		//a := token.Valid
		//b := claims.Valid()
		//fmt.Println("ok")
		return claims, nil
	} else {
		//fmt.Println("failse")
		return nil, err
	}
}