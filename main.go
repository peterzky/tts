package main

import (
	"fmt"
	"os/exec"

	"flag"

	"io/ioutil"
	"os"

	"github.com/peterzky/tts/say"
)

var clip, selec, pipe, debug bool
var input string
var word int

func init() {
	flag.BoolVar(&clip, "clip", false, "send clipboard")
	flag.BoolVar(&pipe, "p", false, "use stdin as input")
	flag.BoolVar(&selec, "sel", false, "send selection")
	flag.BoolVar(&debug, "debug", false, "debug output")
	flag.StringVar(&input, "t", "", "send text")
	flag.IntVar(&word, "w", 10, "word per session")
}

// download text to speech sequentially
func downloader(ch chan say.VoicePart, voiceParts []say.VoicePart) {

	for _, vp := range voiceParts {
		say.Download(vp)
		fmt.Println("download finished %d", vp.Index)
		ch <- vp
	}
}

// play audio with mpg123
func player(ch chan say.VoicePart, counter int, debug bool) {

	for i := 0; i < counter; i++ {
		vp := <-ch
		if debug {
			fmt.Printf("Index: %d\nMessage: %s\nFileName: %s\n",
				vp.Index, vp.Message, vp.FileName)
			fmt.Println("---------------------------")
		}
		tmux := exec.Command("mpg123", vp.FileName)
		tmux.Run()
	}
}

func main() {
	flag.Parse()

	var text string

	if pipe {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		text = string(b)

	}

	if clip {
		xsel := exec.Command("xsel", "-b")
		out, err := xsel.Output()
		text = string(out)
		if err != nil {
			panic(err)
		}
	}

	if selec {
		xsel := exec.Command("xsel", "-o")
		out, err := xsel.Output()
		text = string(out)
		if err != nil {
			panic(err)
		}

	}
	if !clip && !selec && input != "" {
		text = input
	}

	voiceParts := say.Split(text, word)

	download_finish := make(chan say.VoicePart, 30)

	go downloader(download_finish, voiceParts)

	player(download_finish, len(voiceParts), debug)

}
