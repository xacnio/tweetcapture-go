package examples

import (
	"fmt"
	"github.com/Xacnio/tweetcapture-go/pkg/tweetcapture"
)

func TweetScreenshot(url string, nightmode uint8) {
	tss := tweetcapture.NewTweetScreenshot(url)
	tss.SavePath = fmt.Sprintf("./%d.png", nightmode)
	tss.Opts.NightMode = nightmode
	err := tss.Screenshot()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Saved: %s\n", tss.SavePath)
	}
}