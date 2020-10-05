package stockbot

import (
	//"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func AppStockMentionHandler(w http.ResponseWriter, r *http.Request) {
	// get the request body
	defer func() {
		_ = r.Body.Close()
	}()

	body, _ := ioutil.ReadAll(r.Body)

	// unmarshal the whole body (JSON) into a map
	m := make(map[string]interface{})
	err := json.Unmarshal(body, &m)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error unmarshalling body: %v", err)
		return
	}

	// if it's the first them this server is registered with slack,
	// a challenge code needs to be returned
	//	fmt.Fprintf(w, "%s", m["challenge"])
	//}

	// extract the event field into another map
	m1 := m["event"].(map[string]interface{})

	// get the text field
	text := fmt.Sprintf("%v", m1["text"])
	str := strings.Split(text, "<bot user's ID>")
	symbol := strings.Trim(str[1], " ")

	// get the channel id
	channel := fmt.Sprintf("%v", m1["channel"])

	token := "User's Slack Bot Token"

	// get stock price
	price, err := getQuote(symbol)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error calling finnhub API: %v", err)
		return
	}

	// send the price to slack channel
	err = sendMessage(token, channel, price)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error sending message to Slack: %v", err)
		return
	}
}

func sendMessage(token, channel, text string) error {
	postURL := "https://slack.com/api/chat.postMessage"
	data := url.Values{"token": {token}, "channel": {channel}, "text": {text}}
	_, err := http.PostForm(postURL, data)

	return err
}

func getQuote(sym string) (string, error) {
	sym = strings.ToUpper(sym)

	fhUrl := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s", sym)
	resp, err := http.Get(fhUrl)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	m := make(map[string]float32)
	err = json.Unmarshal(body, &m)
	if err != nil {
		return "", err
	}

	fhUrlForName := fmt.Sprintf("https://finnhub.io/api/v1/stock/profile2?symbol=%s", sym)
	resp, err = http.Get(fhUrlForName)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err = ioutil.ReadAll(resp.Body)
	m1 := make(map[string]string)
	err = json.Unmarshal(body, &m1)
	if err != nil {
		return "", err
	}

	var s string

	if len(m) == 1 {
		s = "_" + sym + " is not a valid trading name_"
	} else {
		s = "*" + m1["name"] + " (" + sym + ") " + " Stock Price* \n" +
			"_current_: $" + fmt.Sprintf("%.2f", m["c"]) + "\n" +
			"_high_: $" + fmt.Sprintf("%.2f", m["h"]) + "\n" +
			"_low_: $" + fmt.Sprintf("%.2f", m["l"]) + "\n" +
			"_open_: $" + fmt.Sprintf("%.2f", m["o"]) + "\n" +
			"_previous close_: $" + fmt.Sprintf("%.2f", m["pc"]) + "\n" +
			"_timestamp_: " + time.Now().UTC().String()
	}

	return s, nil
}
