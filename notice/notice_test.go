package notice

import (
	"fmt"
	"testing"
)

func Test_Notice(test *testing.T) {
	notice := NewNotice()
	//配置钉钉 更换为自己的配置信息
	notice.DingConfig.WebhookURL = "https://oapi.dingtalk.com/robot/send?access_token=<DingToken>"

	//配置微信 更换为自己的配置信息
	notice.WxConfig.CorpID = "<CorpID"
	notice.WxConfig.CorpSecret = "<CorpSecret>"

	//配置接收人
	notice.Msg.AgentID = "<AgentID>"
	notice.Msg.User = "<User>"
	var (
		rel string
		err error
	)
	rel, err = notice.Qweixin()
	if err != nil {
		fmt.Printf("Qweixin:%s\n", err)
	} else {
		fmt.Printf("Qweixin发送成功:%s\n", rel)
	}

	rel, err = notice.DingTalk()
	if err != nil {
		fmt.Printf("DingTalk：%s\n", err)
	} else {
		fmt.Printf("DingTalk发送成功:%s\n", rel)
	}

	rel, err = notice.Send(2)
	if err != nil {
		fmt.Printf("Send：%s\n", err)
	} else {
		fmt.Printf("Send发送成功:%s\n", rel)
	}

}
