package main

import (
	"fmt"
	"os/exec"

	"sync"

	"flag"

	"io/ioutil"
	"os"

	"strings"

	"github.com/peterzky/misc/tts/say"
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
	flag.IntVar(&word, "w", 150, "word per session")
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
	if debug {
		for _, v := range voiceParts {
			fmt.Printf("Index: %d\nMessage: %s\nFileName: %s\n", v.Index, v.Message, v.FileName)
			fmt.Println("---------------------------")
		}
	}
	var wg sync.WaitGroup

	for _, vp := range voiceParts {
		wg.Add(1)
		go func(vp say.VoicePart) {
			defer wg.Done()
			say.Download(vp)
		}(vp)
	}
	wg.Wait()

	var files []string
	for _, vp := range voiceParts {
		files = append(files, vp.FileName)
	}
	cmd := fmt.Sprintf("mpg123 %s", strings.Join(files, " "))

	tmux := exec.Command("tmux", "new-window", "-n", "tts", cmd)
	tmux.Run()

}
