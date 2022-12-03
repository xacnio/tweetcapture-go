package tweetcapture

import (
	"errors"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type tweetScreenshotConf struct {
	URL      string
	Opts     tweetScreenshotOpts
	SavePath string
}

type tweetScreenshotOpts struct {
	NightMode        uint8
	SeleniumDebug    bool
	ChromeDriverPath string
}

func NewTweetScreenshot(url string) *tweetScreenshotConf {
	var c tweetScreenshotConf
	c.URL = url
	return &c
}

func NewTweetScreenshotNightMode(url string, nightMode uint8) *tweetScreenshotConf {
	c := NewTweetScreenshot(url)
	c.Opts.NightMode = nightMode
	return c
}

func (c *tweetScreenshotConf) Screenshot() error {
	if !c.validURL() {
		return errors.New(fmt.Sprintf("Invalid Tweet URL: %s", c.URL))
	}

	var chromeDriverPath = chromeDriverPath()
	if len(c.Opts.ChromeDriverPath) > 0 {
		chromeDriverPath = c.Opts.ChromeDriverPath
	}
	const port = 8080

	var opts []selenium.ServiceOption
	selenium.SetDebug(c.Opts.SeleniumDebug)

	if len(c.SavePath) < 1 {
		c.SavePath = "./" + c.defaultFileName()
	}

	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--remote-debugging-port=9222",
			"--headless",
			"--test-type",
			"--disable-logging",
			"--ignore-certificate-errors",
			"--disable-dev-shm-usage",
			"--window-size=768,2000",
		},
	}

	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		return err
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		return err
	}
	defer wd.Quit()

	if err := wd.Get(c.URL); err != nil {
		return err
	}
	wd.AddCookie(c.nightModeCookie())
	wd.Refresh()

	tryCount := 0
	for {
		base := fmt.Sprintf("//a[@href=\"%s\"]/ancestor::article/..", c.baseHref())
		we, err := wd.FindElement(selenium.ByXPATH, base)
		if err == nil {
			time.Sleep(time.Second * 1)
			bs, err := we.Screenshot(true)
			if err == nil {
				os.WriteFile(c.SavePath, bs, 0644)
			}
			return nil
		} else {
			tryCount += 1
			time.Sleep(time.Millisecond * 100)
			if tryCount >= 50*5 {
				return errors.New("timeout: tweet content not found")
			}
		}
	}
}

func (c *tweetScreenshotConf) baseHref() string {
	re := regexp.MustCompile("^https?:\\/\\/([A-Za-z0-9.]+)?twitter\\.com\\/(?:#!\\/)?(\\w+)\\/status(es)?\\/(\\d+)")
	test := re.FindStringSubmatch(c.URL)
	return fmt.Sprintf("/%s/status/%s", strings.ToLower(test[2]), strings.ToLower(test[4]))
}

func (c *tweetScreenshotConf) defaultFileName() string {
	re := regexp.MustCompile("^https?:\\/\\/([A-Za-z0-9.]+)?twitter\\.com\\/(?:#!\\/)?(\\w+)\\/status(es)?\\/(\\d+)")
	test := re.FindStringSubmatch(c.URL)
	return fmt.Sprintf("tweetcapture_@%s_%s.png", strings.ToLower(test[2]), strings.ToLower(test[4]))
}

func (c *tweetScreenshotConf) validURL() bool {
	result, _ := regexp.MatchString("^https?:\\/\\/([A-Za-z0-9.]+)?twitter\\.com\\/(?:#!\\/)?(\\w+)\\/status(es)?\\/(\\d+)", c.URL)
	return result
}

func (c *tweetScreenshotConf) nightModeCookie() *selenium.Cookie {
	timestamp := uint(time.Now().Unix() + 3600)
	nightMode := &selenium.Cookie{Name: "night_mode", Value: fmt.Sprintf("%d", c.Opts.NightMode), Path: "/", Domain: ".twitter.com", Expiry: timestamp, Secure: true}
	return nightMode
}

func chromeDriverPath() string {
	chromeDriverEnv := os.Getenv("CHROME_DRIVER")
	if len(chromeDriverEnv) > 0 {
		return chromeDriverEnv
	}
	if runtime.GOOS == "windows" {
		return "C:/bin/chromedriver.exe"
	}
	return "/usr/local/bin/chromedriver"
}
