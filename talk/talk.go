package talk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/soeyusuke/fb-bot/types"
)

const talkApiUrl = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk" //recruit talk API

func Talk(m string) string {
	params := url.Values{
		"apikey": {os.Getenv("TALKAPIID")},
		"query":  {m},
	}
	json := types.TalkJson{}

	err := post(talkApiUrl, params, &json)
	if err != nil {
		return "ちょっとよくわかりません"
	} else {
		return json.Results[0].Reply
	}
}

func post(url string, params url.Values, out interface{}) error {
	resp, err := http.PostForm(url, params)
	if err != nil {

		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, out)
	if err != nil {
		return err
	}

	return nil
}
