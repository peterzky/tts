package say

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/net/proxy"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

const PROXY_ADDR = "127.0.0.1:1080"
const PROXY = false

func Download(vp VoicePart) {
	var httpclient *http.Client
	if PROXY {
		dialer, _ := proxy.SOCKS5("tcp", PROXY_ADDR, nil, proxy.Direct)
		httpTransport := &http.Transport{}
		httpTransport.Dial = dialer.Dial
		httpclient = &http.Client{Transport: httpTransport}

	} else {
		httpclient = &http.Client{}
	}
	sess, err := session.NewSession(&aws.Config{
		Region:     aws.String("ap-northeast-1"),
		HTTPClient: httpclient,
	})

	svc := polly.New(sess)
	input := &polly.SynthesizeSpeechInput{
		LexiconNames: []*string{},
		OutputFormat: aws.String("mp3"),
		SampleRate:   aws.String("8000"),
		Text:         aws.String(vp.Message),
		TextType:     aws.String("text"),
		VoiceId:      aws.String("Joanna"),
	}

	result, err := svc.SynthesizeSpeech(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case polly.ErrCodeTextLengthExceededException:
				fmt.Println(polly.ErrCodeTextLengthExceededException, aerr.Error())
			case polly.ErrCodeInvalidSampleRateException:
				fmt.Println(polly.ErrCodeInvalidSampleRateException, aerr.Error())
			case polly.ErrCodeInvalidSsmlException:
				fmt.Println(polly.ErrCodeInvalidSsmlException, aerr.Error())
			case polly.ErrCodeLexiconNotFoundException:
				fmt.Println(polly.ErrCodeLexiconNotFoundException, aerr.Error())
			case polly.ErrCodeServiceFailureException:
				fmt.Println(polly.ErrCodeServiceFailureException, aerr.Error())
			case polly.ErrCodeMarksNotSupportedForFormatException:
				fmt.Println(polly.ErrCodeMarksNotSupportedForFormatException, aerr.Error())
			case polly.ErrCodeSsmlMarksNotSupportedForTextTypeException:
				fmt.Println(polly.ErrCodeSsmlMarksNotSupportedForTextTypeException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	file, _ := os.Create(vp.FileName)
	defer file.Close()

	io.Copy(file, result.AudioStream)

}
