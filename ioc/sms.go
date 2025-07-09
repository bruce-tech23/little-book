package ioc

import (
	"geektime-basic-learning2/little-book/internal/service/sms"
	"geektime-basic-learning2/little-book/internal/service/sms/localsms"
	"geektime-basic-learning2/little-book/internal/service/sms/tencent"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSMSService() sms.Service {
	return localsms.NewService()
	// 如果有需要，就可以切换用腾讯的短信服务
	// return initTencentSMSService()
}

func initTencentSMSService() sms.Service {
	secretId, ok := os.LookupEnv("自己命名的环境变量SecretId")
	if !ok {
		panic("找不到腾讯 SMS 的 secret id")
	}
	secretKey, ok := os.LookupEnv("自己命名的环境变量SecretKey")
	if !ok {
		panic("找不到腾讯 SMS 的 secret key")
	}
	c, err := tencentSMS.NewClient(common.NewCredential(secretId, secretKey), "ap-nanjing", profile.NewClientProfile())
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "xxx", "aa")
}
