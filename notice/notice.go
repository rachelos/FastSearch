package notice

import (
	"fmt"

	"gitee.com/rachel_os/fastsearch/notice/dingding"
	"gitee.com/rachel_os/fastsearch/notice/qyweixin"
)

type Message struct {
	AgentID string `json:"agent_id"` // 接收人ID
	User    string `json:"user"`     // 接收人
	MsgType string `json:"msgtype"`  // 消息类型，这里使用text
	Text    struct {
		Content string `json:"content"` // 消息内容
	} `json:"text"`
	// 如果需要发送Markdown或其他类型消息，可以添加相应字段
}
type Config struct {
	CorpID     string `json:"corp_id"`
	CorpSecret string `json:"corp_secret"`
}
type DingConfig struct {
	WebhookURL string
}
type SendType int

const (
	QyWeixin SendType = iota + 1
	DingTalk
)

type Notice struct {
	Msg        *Message
	SendType   SendType
	WxConfig   *Config
	DingConfig *DingConfig
}

func NewNotice() *Notice {
	return &Notice{
		Msg:        &Message{},
		WxConfig:   &Config{},
		DingConfig: &DingConfig{},
	}
}
func (notice *Notice) Qweixin() (string, error) {
	corpID := notice.WxConfig.CorpID
	corpSecret := notice.WxConfig.CorpSecret
	agentID := notice.Msg.AgentID
	user := notice.Msg.User
	msg := notice.Msg.Text.Content

	accessToken, err := qyweixin.GetAccessToken(corpID, corpSecret)
	if err != nil {
		return "", err
	}
	body, err := qyweixin.SendAppMessage(accessToken, agentID, user, msg)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return "", err
	}
	return body, nil
}

func (notice *Notice) DingTalk() (string, error) {
	// "https://oapi.dingtalk.com/robot/send?access_token=YOUR_ACCESS_TOKEN"
	webhookURL := notice.DingConfig.WebhookURL
	message := dingding.DingTalkMessage{
		MsgType: notice.Msg.MsgType,
		Text: struct {
			Content string `json:"content"`
		}{
			Content: notice.Msg.Text.Content,
		},
	}
	rel, err := dingding.SendDingTalkMessage(webhookURL, message)
	return rel, err
}

func (notice *Notice) Send(sendType SendType) (string, error) {
	notice.SendType = sendType
	switch notice.SendType {
	case QyWeixin:
		return notice.Qweixin()
	case DingTalk:
		return notice.DingTalk()
	default:
	}
	return "", nil
}
