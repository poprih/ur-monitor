package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ResponseItem represents the structure of the response data
type ResponseItem struct {
	Name           string `json:"name"`
	Floortype      string `json:"floor"`
	RentNormal     string `json:"rent_normal"`
	RoomDetailLink string `json:"roomDetailLink"`
}

// sendLineMessage sends a LINE Notify message
func sendLineMessage(token, message string) error {
	apiURL := "https://notify-api.line.me/api/notify"
	data := url.Values{}
	data.Set("message", message)

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send LINE message, status: %d", resp.StatusCode)
	}
	return nil
}

// fetchData requests UR property data
func fetchData(danchi string) ([]ResponseItem, error) {
	url := "https://chintai.r6.ur-net.go.jp/chintai/api/bukken/detail/detail_bukken_room/"
	postData := fmt.Sprintf("rent_low=&rent_high=&floorspace_low=&floorspace_high=&shisya=80&danchi=%s&shikibetu=0&newBukkenRoom=&orderByField=0&orderBySort=0&pageIndex=0&pageIndex=0&sp=", danchi)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postData)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []ResponseItem
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// 补全 `roomDetailLink` 的域名
	for i := range data {
		data[i].RoomDetailLink = "https://www.ur-net.go.jp" + data[i].RoomDetailLink
	}

	return data, nil
}

// StartMonitorJob starts the monitoring job
func StartMonitorJob(danchis []string, lineToken string) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		// 1~5 分钟随机等待
		delay := time.Duration(rand.Intn(5)+1) * time.Minute
		fmt.Printf("Waiting for %v before sending request...\n", delay)
		time.Sleep(delay)

		var allResults []ResponseItem
		for _, danchi := range danchis {
			fmt.Println("Fetching data for danchi:", danchi)
			data, err := fetchData(danchi)
			if err != nil {
				fmt.Println("Request failed:", err)
				continue
			}
			allResults = append(allResults, data...)
		}

		// 整合所有数据后发送
		if len(allResults) > 0 {
			var messageBuilder strings.Builder
			for _, item := range allResults {
				messageBuilder.WriteString(fmt.Sprintf("🏠 Name: %s\n🏢 Floor: %s\n💰 Rent: %s\n🔗 Link: %s\n\n",
					item.Name, item.Floortype, item.RentNormal, item.RoomDetailLink))
			}
			message := messageBuilder.String()
			err := sendLineMessage(lineToken, message)
			if err != nil {
				fmt.Println("Failed to send LINE message:", err)
			} else {
				fmt.Println("LINE message sent successfully!")
			}
		} else {
			fmt.Println("No data found, skipping LINE notification.")
		}
	}
}
