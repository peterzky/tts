package say

import (
	"fmt"
	"strings"
)

type VoicePart struct {
	Message  string
	FileName string
	Index    int
}

type ByIndex []VoicePart

// implement sort interface
func (b ByIndex) Len() int {
	return len(b)
}

func (b ByIndex) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByIndex) Less(i, j int) bool {
	return b[i].Index < b[j].Index
}

func getName(k int) string {
	return fmt.Sprintf("/tmp/tts_%d.mp3", k)
}

func Split(longStr string, length int) []VoicePart {
	var list []VoicePart
	words := strings.Split(longStr, " ")
	var k int = 1
	for i := 0; i < len(words); i += length {
		if len(words) < length {
			str := longStr
			list = append(list, VoicePart{str, getName(k), k})
		}

		if i%length == 0 && i != 0 {
			str := strings.Join(words[i-length:i], " ")
			list = append(list, VoicePart{str, getName(k), k})
			k++

		}
		if len(words)-i < length && i != 0 {
			str := strings.Join(words[i:], " ")
			list = append(list, VoicePart{str, getName(k), k})

		}
	}

	return list
}
