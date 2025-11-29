package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	chromeDriverPath = "/opt/homebrew/bin/chromedriver"
	chromePort       = 9515
	baseURL          = "http://localhost:8081"
)

func main() {
	service, err := selenium.NewChromeDriverService(chromeDriverPath, chromePort)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	chromeCaps := chrome.Capabilities{
    	Args: []string{
        	//"--headless=new",
        	"--disable-gpu",
        	"--no-sandbox",
        	"--disable-dev-shm-usage",
        	"--remote-allow-origins=*",
			"--host-resolver-rules=MAP localhost 127.0.0.1",
    	},
	}

	caps := selenium.Capabilities{}
	caps.AddChrome(chromeCaps)

	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", chromePort))
	if err != nil {
		panic(err)
	}
	defer driver.Quit()

	fmt.Println("Starting Selenium backend test suite")

	driver.Get(baseURL)

	userPayload := `{"username":"testuser","email":"testuser@example.com","password":"pass123"}`
	createUserResp := fetchPost(driver, "/v1/create_user", userPayload)
	fmt.Println("CreateUser response:", createUserResp)

	loginPayload := `{"email":"testuser@example.com","password":"pass123"}`
	loginResp := fetchPost(driver, "/v1/login", loginPayload)
	fmt.Println("Login response:", loginResp)

	accessToken := extractToken(loginResp, "access_token")
	if accessToken == "" {
		panic("Login failed : access_token missing")
	}
	fmt.Println("Login successful")

	marketPayload := `{"name":"BTC-USD","base_currency":"BTC","quote_currency":"USD"}`
	marketResp := fetchPost(driver, "/v1/markets", marketPayload)
	fmt.Println("CreateMarket response:", marketResp)
	marketID := extractID(marketResp, "market_id")

	orderPayload := fmt.Sprintf(`{"market_id":"%s","side":"buy","quantity":1.0,"price":30000}`, marketID)
	orderResp := fetchPost(driver, "/v1/orders", orderPayload)
	fmt.Println("CreateOrder response:", orderResp)
	orderID := extractID(orderResp, "order_id")

	walletPayload := `{"currency":"BTC"}`
	walletResp := fetchPost(driver, "/v1/wallets", walletPayload)
	fmt.Println("CreateWallet response:", walletResp)
	walletID := extractID(walletResp, "wallet_id")

	txPayload := fmt.Sprintf(`{"wallet_id":"%s","amount":0.5,"type":"deposit"}`, walletID)
	txResp := fetchPost(driver, "/v1/transactions", txPayload)
	fmt.Println("CreateTransaction response:", txResp)

	tradePayload := fmt.Sprintf(`{"market_id":"%s","order_id":"%s","quantity":1.0,"price":30000}`, marketID, orderID)
	tradeResp := fetchPost(driver, "/v1/trades", tradePayload)
	fmt.Println("CreateTrade response:", tradeResp)

	driver.Get(baseURL + "/oauth/google/login")
	fmt.Println("OAuth Google Login page opened (check browser if not headless)")

	fmt.Println("Selenium backend test suite finished")
}

func fetchPost(driver selenium.WebDriver, path, payload string) string {
    script := fmt.Sprintf(`
        const payload = %s;
        return fetch("%s%s", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
            mode: "cors"
        })
        .then(async r => {
            const text = await r.text();
            return JSON.stringify({ ok: r.ok, status: r.status, statusText: r.statusText, body: text });
        })
        .catch(e => {
            return JSON.stringify({
                error: e.toString(),
                name: e.name,
                message: e.message,
                stack: e.stack
            });
        });
    `, payload, baseURL, path)

    res, err := driver.ExecuteScript(script, nil)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to execute fetch POST script")
    }
    fmt.Printf("DEBUG fetchPost response for %s: %v\n", path, res)
    return res.(string)
}



func extractToken(jsonStr, key string) string {
	var data map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &data)
	if val, ok := data[key]; ok {
		return val.(string)
	}
	return ""
}

func extractID(jsonStr, key string) string {
	var data map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &data)
	if val, ok := data[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}
