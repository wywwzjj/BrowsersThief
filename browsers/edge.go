package browsers

import "os"

type Edge struct {
	*Chromium
}

func NewEdge() *Edge {
	edgeUserDataPath := os.Getenv("LOCALAPPDATA") + "\\MicroSoft\\Edge\\User Data"
	edge := &Edge{newChromium(edgeUserDataPath)}

	return edge
}
