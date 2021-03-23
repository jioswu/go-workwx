package workwx

import (
	"encoding/json"
	"fmt"
)

//send_welcome_msg 发送新客户欢迎语
type ReqSendWelcomeMsg struct {
	WelcomeCode string             `json:"welcome_code"`
	Text        SendWelcomeMsgText `json:"text,omitempty"`
	//Image       SendWelcomeMsgImage       `json:"image,omitempty"`
	//Link        SendWelcomeMsgLink        `json:"link,omitempty"`
	Miniprogram SendWelcomeMsgMiniprogram `json:"miniprogram,omitempty"`
}

type SendWelcomeMsgText struct {
	Content string `json:"content"`
}

type SendWelcomeMsgImage struct {
	MediaId string `json:"media_id"` //图片的media_id，可以通过素材管理接口获得
	PicUrl  string `json:"pic_url"`  //图片的链接，仅可使用上传图片接口得到的链接
}

type SendWelcomeMsgLink struct {
	Title  string `json:"title"`  //图文消息标题，最长为128字节
	Picurl string `json:"picurl"` //图文消息封面的url
	Desc   string `json:"desc"`   //图文消息的描述，最长为512字节
	Url    string `json:"url"`    //图文消息的链接
}

type SendWelcomeMsgMiniprogram struct {
	// Appid 小程序appid，必须是有在本企业安装授权的小程序，否则会被忽略
	Appid string `json:"appid"`
	// page 小程序的页面路径
	Page string `json:"page"`
	// Title 企业对外简称，需从已认证的企业简称中选填。可在“我的企业”页中查看企业简称认证状态。
	Title      string `json:"title"`
	PicMediaId string `json:"pic_media_id"` //小程序消息封面的mediaid，封面图建议尺寸为520*416
}

func (x ReqSendWelcomeMsg) intoBody() ([]byte, error) {
	result, err := json.Marshal(x)
	fmt.Printf("ReqSendWelcomeMsg initToBody:%s\n", string(result))
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

type RespSendWelcomeMsg struct {
	respCommon
}

// execExternalContactSendWelcomeMsg 发送新客户欢迎语
func (c *WorkwxApp) execExternalContactSendWelcomeMsg(req ReqSendWelcomeMsg) (resp RespSendWelcomeMsg, err error) {
	resp = RespSendWelcomeMsg{}
	err = c.executeQyapiJSONPost("/cgi-bin/externalcontact/send_welcome_msg", req, &resp, true)
	if err != nil {
		return
	}
	if bizErr := resp.TryIntoErr(); bizErr != nil {
		return
	}

	return resp, nil
}

func (c *WorkwxApp) ExecExternalContactSendWelcomeMsg(req *ReqSendWelcomeMsg) (resp RespSendWelcomeMsg, err error) {
	return c.execExternalContactSendWelcomeMsg(*req)
}
