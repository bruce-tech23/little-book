package tencent

import (
	"context"
	"errors"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
) // 引入sms

type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func (s Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = common.StringPtr(tplId)
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(numbers)
	response, err := s.client.SendSms(request)
	// 处理异常
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		// 因为 status 是指针，所以要判断是否是 nil
		if statusPtr != nil {
			// 基本不可能进来这里
			continue
		}
		status := *statusPtr
		// 发送失败。注意这里在实际业务中批量发送时要区分出发送成功与失败的，失败的业务需求可能需要再次发送。
		if status.Code == nil {
			return errors.New("message sent failed")
		} else if *(status.Code) != "Ok" {
			return fmt.Errorf("message sent failed. code: %s, msg: %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func NewService(client *sms.Client, appId string, signature string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signature,
	}
}
