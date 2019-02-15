package say

import (
	"fmt"
	"regexp"
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

func SplitSentence(longStr string) []string {
	var sentenceList []string
	var sentence string
	remove_new_line := strings.NewReplacer("\n", " ", "\t", " ")
	str := remove_new_line.Replace(longStr)

	for _, ch := range str {
		char := string(ch)
		end, _ := regexp.Match("[.!?]", []byte(char))
		middle, _ := regexp.Match(("[,.?!]"), []byte(char))
		sentenceLen := len(sentence)

		switch {
		case end && sentenceLen > 50:
			sentence += char
			sentenceList = append(sentenceList, sentence)
			sentence = ""
		case middle && sentenceLen > 200:
			sentence += char
			sentenceList = append(sentenceList, sentence)
			sentence = ""
		default:
			sentence += char
		}

	}

	if len(sentence) != 0 {
		sentenceList = append(sentenceList, sentence)
	}
	return sentenceList

}

// split string into sentences
func Split(longStr string, length int) []VoicePart {
	var list []VoicePart

	for i, s := range SplitSentence(longStr) {
		list = append(list, VoicePart{s, getName(i), i})
	}

	return list
}
