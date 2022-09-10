// Package internal
// Time    : 2022/9/7 22:01
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package internal

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
	"os/exec"
	"time"
)

type OaWeb struct {
	IsHeadless  bool
	BrowserPath string
	Logger      *log.Entry
	Browser     *rod.Browser
	Launcher    *launcher.Launcher
	Page        *rod.Page
}

type captchaReq struct {
	Base64Str string `json:"base64str"`
}

type captchaRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func NewOaWebLogin(ctx context.Context, headless bool, logger *log.Entry) (*OaWeb, error) {
	o := NewOaWeb(headless, logger)
	err := o.Start()
	if err != nil {
		return nil, err
	}
	err = o.GoLoginPage(ctx)
	if err != nil {
		return nil, err
	}
	err = o.LoginOa(ctx)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second * 1)
	return o, nil
}

func NewOaWeb(headless bool, logger *log.Entry) *OaWeb {
	return &OaWeb{
		IsHeadless: headless,
		Logger:     logger,
	}
}

func (o *OaWeb) Start() error {
	err := o.FindDefaultBrowserPath()
	o.Logger.Infof("browser path: %s", o.BrowserPath)
	o.Launcher = launcher.New().Bin(o.BrowserPath)
	o.Launcher.Logger(o.Logger.Writer())
	if o.IsHeadless {
		o.Launcher.Headless(true)
		o.Launcher.Set("disable-gpu", "true")
	} else {
		o.Launcher.Headless(false)
	}
	o.Launcher.Set("autoplay-policy", "no-user-gesture-required").Set("mute-audio")

	u, err := o.Launcher.Launch()
	if err != nil {
		return err
	}
	o.Browser = rod.New().ControlURL(u)
	err = o.Browser.Connect()

	if err != nil {
		return err
	}
	return nil
}

func (o *OaWeb) Stop() error {
	_ = o.Browser.Close()
	return nil
}

func (o *OaWeb) LogErr(err error, message string, args ...interface{}) error {
	err = errors.Wrapf(err, message, args...)
	o.Logger.Error(err)
	return err
}

func (o *OaWeb) FindDefaultBrowserPath() error {
	chromePath := viper.GetString("chrome.path")

	if len(chromePath) > 0 {
		validPath, err := exec.LookPath(chromePath)
		if err != nil {
			o.Logger.Errorf("Could not find specified chrome path: %s, err: %v", chromePath, err)

		} else {
			o.BrowserPath = validPath
			return nil
		}
	}
	path, isFound := launcher.LookPath()
	if isFound {
		o.BrowserPath = path
		viper.Set("chrome.path", o.BrowserPath)
		return nil
	} else {
		return o.LogErr(fmt.Errorf("can not find default browser path"), "")
	}
}

func (o *OaWeb) GoLoginPage(ctx context.Context) error {
	wrPageUrl := viper.GetString("oa.login_url")
	if len(wrPageUrl) == 0 {
		return o.LogErr(nil, "oa.login_url in .ndr.toml empty")
	}
	u, err := url.Parse(wrPageUrl)
	if err != nil {
		return o.LogErr(err, "oa.login_url in .ndr.toml")
	}

	o.Page, err = o.Browser.Page(proto.TargetCreateTarget{URL: u.String()})
	if err != nil {
		return o.LogErr(err, "oa.login_url in .ndr.toml invalid")
	}
	o.Page.Context(ctx)
	err = o.Page.WaitLoad()
	if err != nil {
		return o.LogErr(err, "load oa.login_url err")
	}
	return nil
}

func (o *OaWeb) HasX(selector string) (*rod.Element, error) {
	isFind, element, err := o.Page.HasX(selector)
	if !isFind {
		if err != nil {
			return nil, o.LogErr(err, "find element %s err", selector)
		} else {
			return nil, o.LogErr(nil, "not find element %s", selector)
		}
	}
	return element, nil
}

func (o *OaWeb) InputTextX(selector, input string) error {
	element, err := o.HasX(selector)
	if err != nil {
		return err
	}
	err = rod.Try(func() {
		element.MustSelectAllText().MustInput("")
	})
	if err != nil {
		return o.LogErr(err, "input element : %s, clear text", selector)
	}

	err = element.Input(input)
	if err != nil {
		return o.LogErr(err, "input element : %s, text: %s", selector, input)
	}
	return nil
}

func (o *OaWeb) ClickBtnX(selector string) error {
	element, err := o.HasX(selector)
	if err != nil {
		return err
	}
	err = element.Click(proto.InputMouseButtonLeft)
	if err != nil {
		return o.LogErr(err, "click %s failed", selector)
	}
	return nil
}

func (o *OaWeb) ElementAttribute(element *rod.Element, name string) (value *string, err error) {
	v, err := element.Attribute(name)
	if err != nil {
		return nil, o.LogErr(err, "no attribute: %s", name)
	}
	return v, nil
}

func (o *OaWeb) GetCaptchaPath(selector string) (*url.URL, error) {
	element, err := o.HasX(selector)
	if err != nil {
		return nil, err
	}
	v, err := o.ElementAttribute(element, "src")
	if err != nil {
		return nil, o.LogErr(err, "selector: %s", selector)
	}
	uStr := viper.GetString("oa.captcha_host") + *v
	u, err := url.Parse(uStr)
	if err != nil {
		return nil, o.LogErr(err, "parse captcha url: %s", uStr)
	}
	return u, nil
}

func (o *OaWeb) GetCaptchaStr(u *url.URL) (string, error) {
	bin, err := o.Page.GetResource(u.String())
	if err != nil {
		return "", o.LogErr(err, "GetCaptcha: %s", u.String())
	}
	b64Str := base64.StdEncoding.EncodeToString(bin)
	o.Logger.Debugf("captcha jpg base64 encoded: %s", b64Str)

	capReq := &captchaReq{Base64Str: b64Str}
	capRes := &captchaRes{}

	client := req.C()
	res, err := client.R().SetBody(
		capReq).SetResult(capRes).Post(viper.GetString("captcha.url"))
	if err != nil {
		return "", o.LogErr(err, "parse captcha base64 %s, failed", b64Str)
	}
	if res.IsSuccess() {
		if capRes.Code == 0 {
			return capRes.Data, nil
		} else {
			return "", o.LogErr(fmt.Errorf("parse captcha base64 %s, error: %s", b64Str, capRes.Msg), "")
		}
	} else {
		return "", o.LogErr(fmt.Errorf("parse captcha base64 %s, http error: %v", b64Str, res.Error()), "")
	}
}

func (o *OaWeb) RetryLoginBtn(u *url.URL, retryCnt int) error {
	var e error = nil
	var rro *proto.RuntimeRemoteObject
	for i := 0; i < retryCnt; i++ {
		captcha, err := o.GetCaptchaStr(u)
		if err != nil {
			e = err
		} else {
			err = o.InputTextX(`//input[@name="vcode"]`, captcha)
			if err != nil {
				e = errors.Wrapf(err, "input captcha: %s failed", captcha)
			} else {
				err = o.ClickBtnX(`//button[@type="button" and string()="立即登录"]`)
				if err != nil {
					e = errors.Wrap(err, "click login button failed")
				} else {
					time.Sleep(time.Second * 2)
					rro, err = o.Page.Eval(`() => window.location.host`)
					if err != nil {
						e = o.LogErr(err, "run js () => {return window.location.host} failed")
					} else {
						if rro.Value.Str() == "oa.jss.com.cn" {
							err = nil
							break
						}
					}
					// e, err := o.HasX(`label[@for="vcode" and @style="display: inline;"]`)
					// if err != nil {
					// 	if strings.Contains(fmt.Sprintln(err), "not find")
					// } else {
					// 	continue
					// }
				}
			}
		}
		//
	}
	if e != nil {
		return e
	} else {
		return nil
	}
}

func (o *OaWeb) StuffLoginInfo(ctx context.Context) error {
	err := o.InputTextX("//*[@id=\"usernameInput\"]", viper.GetString("oa.account"))
	if err != nil {
		err = errors.Wrap(err, "input username failed")
		return err
	}
	err = o.InputTextX("//form[1]/div[2]/input", viper.GetString("oa.password"))
	if err != nil {
		err = errors.Wrap(err, "input password failed")
		return err
	}
	return nil
}

func (o *OaWeb) LoginOa(ctx context.Context) error {
	err := o.StuffLoginInfo(ctx)
	if err != nil {
		return err
	}
	u, err := o.GetCaptchaPath("//img[@onclick]")
	if err != nil {
		return err
	}
	err = o.RetryLoginBtn(u, 10)
	if err != nil {
		return err
	}
	// TODO: get captcha string
	return nil
}
