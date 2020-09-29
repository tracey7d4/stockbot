# Simple Slack bot in Go using event API - II
## Stock bot

We will discuss about `stockbot` that returns stock value of a given company at current time.

![stockBot image](docs/image/stockbot.png)  

 For one who comes to this repository first, I have another repository call [`weatherbot`](https://github.com/Tracey7d4/weatherbot) 
 that step by step describe how to build a `weatherbot` that return weather condition of a given city.
Here I only discuss some differences regarding the Go function.

   Similar to `weatherbot`,   we now create Slack bot named `stockbot` and our Go function named `AppStockMentionHandler`.
  Most of steps in this function are similar to the ones in previous function for `weatherbot`,
   the only difference is that we use `getQuote` function here  instead of `getWeather` function. 

```go
    func getQuote(sym string) (string, error) {
    	sym = strings.ToUpper(sym)
    
    	fhUrl := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s", sym)
    	resp, err := http.Get(fhUrl)
    	if err != nil {
    		return "", err
    	}
    
    	defer func(){
    		_ = resp.Body.Close()
    	}()
    
    	body, err := ioutil.ReadAll(resp.Body)
    	m := make(map[string]float32)
    	err = json.Unmarshal(body, &m)
    	if err != nil {
    		return "", err
    	}
    
    	url1 := fmt.Sprintf("https://finnhub.io/api/v1/stock/profile2?symbol=%s", sym)
    	resp1, err1 := http.Get(url1)
    	if err1 != nil {
    		return "", err1
    	}
    
    	defer func(){
    		_ = resp1.Body.Close()
    	}()
    
    	body1, err1 := ioutil.ReadAll(resp1.Body)
    	m1 := make(map[string]string)
    	err = json.Unmarshal(body1, &m1)
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
```
  
  I use `http://finnhub.io API` this time for getting stock value of a company by its ticker symbol. You can find
  its API documentation in [here](https://finnhub.io/docs/api).
  
  In case you want to display the company's trading name along with its ticker symbol,
  you can extract this information from following URL (`url1` variable in the code)
  
   ```shell script
   https://finnhub.io/api/v1/stock/profile2?symbol=<ticker symbol>
   ```
   
   Now let's deploy our `Go` function
   
   ```shell script
   gcloud functions deploy AppStockMentionHandler --runtime go111 --trigger-http
   ```
   When you see that your function has been successfully deployed, go to a Slack Channel and call your bot, remember to
   mention its name.
   
   ```shell script
   @stockbot aapl
   ```  
   ![stockBot image](docs/image/stockbot.png)
   
   
   Hope you enjoy your bots.

### API reference
* [Stock API Documentation](https://finnhub.io/docs/api)
* [Weather bot repository](https://github.com/Tracey7d4/weatherbot)