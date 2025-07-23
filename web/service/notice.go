package service

import (
	"errors"

	"gitee.com/rachel_os/fastsearch/global"
	"gitee.com/rachel_os/fastsearch/notice"
)

type Notice struct {
	Text     string `json:"text"`
	ToUser   string `json:"to_user"`
	ToUserID string `json:"to_user_id"`
}

func NewNotice() *Notice {
	return &Notice{}
}

// 消息通知
func (noc *Notice) Notice(w Notice) ([]string, error) {
	notice := notice.NewNotice()
	notice.Msg.MsgType = "text"
	notice.Msg.AgentID = w.ToUserID
	notice.Msg.User = w.ToUser
	notice.Msg.Text.Content = w.Text

	// 微信配置
	notice.WxConfig.CorpID = global.CONFIG.Notice.QyWeixin.CorpID
	notice.WxConfig.CorpSecret = global.CONFIG.Notice.QyWeixin.CorpSecret

	// 钉钉配置
	notice.DingConfig.WebhookURL = global.CONFIG.Notice.DingTalk.WebHook

	var (
		rel []string
	)
	if !global.CONFIG.Notice.Enable {
		return nil, errors.New("未开启消息通知")
	}
	if global.CONFIG.Notice.DingTalk.Enable {
		v, _ := notice.DingTalk()
		rel = append(rel, v)
	}
	if global.CONFIG.Notice.QyWeixin.Enable {
		v, _ := notice.Qweixin()
		rel = append(rel, v)
	}
	return rel, nil
}
