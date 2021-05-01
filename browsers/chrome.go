package browsers

import (
	"os"
)

type Chrome struct {
	*Chromium
}

func NewChrome() *Chrome {
	chromeUserDataPath := os.Getenv("LOCALAPPDATA") + "\\Google\\Chrome\\User Data"
	chrome := &Chrome{newChromium(chromeUserDataPath)}

	return chrome
}
