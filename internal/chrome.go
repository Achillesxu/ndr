// Package internal
// Time    : 2022/9/7 22:01
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package internal

import (
	"context"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
	"os/exec"
	"time"
)

type Chrome struct {
	IsHeadless  bool
	BrowserPath string
	Logger      *log.Entry
	Browser     *rod.Browser
	Launcher    *launcher.Launcher
	Page        *rod.Page
}

func NewChromeLogin(ctx context.Context, headless bool, logger *log.Entry) (*Chrome, error) {
	c := NewChrome(headless, logger)
	err := c.Start()
	if err != nil {
		return nil, err
	}
	err = c.GoLoginPage(ctx)
	if err != nil {
		return nil, err
	}
	err = c.LoginOa(ctx)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second * 10)
	return c, nil
}

func NewChrome(headless bool, logger *log.Entry) *Chrome {
	return &Chrome{
		IsHeadless: headless,
		Logger:     logger,
	}
}

func (c *Chrome) Start() error {
	err := c.FindDefaultBrowserPath()
	c.Logger.Infof("browser path: %s", c.BrowserPath)
	c.Launcher = launcher.New().Bin(c.BrowserPath)
	c.Launcher.Logger(c.Logger.Writer())
	if c.IsHeadless {
		c.Launcher.Headless(true)
		c.Launcher.Set("disable-gpu", "true")
	} else {
		c.Launcher.Headless(false)
	}
	c.Launcher.Set("autoplay-policy", "no-user-gesture-required").Set("mute-audio")

	u, err := c.Launcher.Launch()
	if err != nil {
		return err
	}
	c.Browser = rod.New().ControlURL(u)
	err = c.Browser.Connect()

	if err != nil {
		return err
	}
	return nil
}

func (c *Chrome) Stop() error {
	_ = c.Browser.Close()
	return nil
}

func (c *Chrome) LogErr(err error, message string, args ...interface{}) error {
	err = errors.Wrapf(err, message, args...)
	c.Logger.Error(err)
	return err
}

func (c *Chrome) FindDefaultBrowserPath() error {
	chromePath := viper.GetString("chrome.path")

	if len(chromePath) > 0 {
		validPath, err := exec.LookPath(chromePath)
		if err != nil {
			c.Logger.Errorf("Could not find specified chrome path: %s, err: %v", chromePath, err)

		} else {
			c.BrowserPath = validPath
			return nil
		}
	}
	path, isFound := launcher.LookPath()
	if isFound {
		c.BrowserPath = path
		viper.Set("chrome.path", c.BrowserPath)
		return nil
	} else {
		return c.LogErr(fmt.Errorf("can not find default browser path"), "")
	}
}

func (c *Chrome) GoLoginPage(ctx context.Context) error {
	wrPageUrl := viper.GetString("oa.login_url")
	if len(wrPageUrl) == 0 {
		return c.LogErr(nil, "oa.login_url in .ndr.toml empty")
	}
	u, err := url.Parse(wrPageUrl)
	if err != nil {
		return c.LogErr(err, "oa.login_url in .ndr.toml")
	}

	c.Page, err = c.Browser.Page(proto.TargetCreateTarget{URL: u.String()})
	if err != nil {
		return c.LogErr(err, "oa.login_url in .ndr.toml invalid")
	}
	c.Page.Context(ctx)
	err = c.Page.WaitLoad()
	if err != nil {
		return c.LogErr(err, "load oa.login_url err")
	}
	return nil
}

func (c *Chrome) HasX(selector string) (*rod.Element, error) {
	isFind, element, err := c.Page.HasX(selector)
	if !isFind {
		if err != nil {
			return nil, c.LogErr(err, "find element err")
		} else {
			return nil, c.LogErr(nil, "not found element via xpath: %s", selector)
		}
	}
	return element, nil
}

func (c *Chrome) InputTextX(selector, input string) error {
	element, err := c.HasX(selector)
	if err != nil {
		return err
	}
	err = element.Input(input)
	if err != nil {
		return errors.Wrapf(err, "input element : %s, text: %s", selector, input)
	}
	return nil
}

func (c *Chrome) ElementAttribute(element *rod.Element, name string) (value *string, err error) {
	v, err := element.Attribute(name)
	if err != nil {
		return nil, c.LogErr(err, "no attribute: %s", name)
	}
	return v, nil
}

func (c *Chrome) GetCaptchaPath(selector string) (*url.URL, error) {
	element, err := c.HasX(selector)
	if err != nil {
		return nil, err
	}
	v, err := c.ElementAttribute(element, "src")
	if err != nil {
		return nil, c.LogErr(err, "selector: %s", selector)
	}
	uStr := viper.GetString("oa.captcha_host") + *v
	u, err := url.Parse(uStr)
	if err != nil {
		return nil, c.LogErr(err, "parse captcha url: %s", uStr)
	}
	return u, nil
}

func (c *Chrome) StuffLoginInfo(ctx context.Context) error {
	err := c.InputTextX("//*[@id=\"usernameInput\"]", viper.GetString("oa.account"))
	if err != nil {
		err = errors.Wrap(err, "input username failed")
		return err
	}
	err = c.InputTextX("//form[1]/div[2]/input", viper.GetString("oa.password"))
	if err != nil {
		err = errors.Wrap(err, "input password failed")
		return err
	}
	return nil
}

func (c *Chrome) LoginOa(ctx context.Context) error {
	err := c.StuffLoginInfo(ctx)
	if err != nil {
		return err
	}
	u, err := c.GetCaptchaPath("//img[@onclick]")
	if err != nil {
		return err
	}
	bin, err := c.Page.GetResource(u.String())
	if err != nil {
		return c.LogErr(err, "GetCaptcha: %s", u.String())
	}
	b64Str := Bin2Base64Str(bin)
	c.Logger.Info(b64Str)
	// TODO: get captcha string
	return nil
}
