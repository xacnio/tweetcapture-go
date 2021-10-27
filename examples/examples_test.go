package examples

import (
	"testing"
)

func TestTweetScreenshot(t *testing.T) {
	TweetScreenshot("https://twitter.com/jack/status/20", 0)
	TweetScreenshot("https://twitter.com/jack/status/20", 1)
	TweetScreenshot("https://twitter.com/jack/status/20", 2)
}