package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/derfinlay/ddns/config"
)

type Record struct {
	Id      string `json:"id"`
	Comment string `json:"comment"`
	Content string `json:"content"`
	Name    string `json:"name"`
}

type ZoneRecordsResult struct {
	Result []Record `json:"result"`
}

var CLOUDFLARE_API_BASE_URL string = "https://api.cloudflare.com/client/v4/"
var IP_API_ENDPOINT string = "https://cloudflare.com/cdn-cgi/trace"

func main() {
	for {
		c, err := config.LoadConfig()

		if err != nil {
			log.Print("Error loading config - retrying in 10 seconds")
			time.Sleep(10 * time.Second)
			continue
		}

		printConig(c)

		log.Print("Starting Cloudflare Record updates...", time.Now())
		go run(c)

		time.Sleep(time.Duration(c.UpdateInterval) * time.Second)
	}
}

func updateRecord(zoneId string, recordId string, newValue string, API_KEY string) error {
	body := []byte(`{"content": "` + newValue + `"}`)
	makePatch(CLOUDFLARE_API_BASE_URL+"/zones/"+zoneId+"/dns_records/"+recordId, API_KEY, body)
	return nil
}

func getRecordsByZoneId(zoneId string, API_KEY string) ([]Record, error) {
	res := makeRequest(CLOUDFLARE_API_BASE_URL+"/zones/"+zoneId+"/dns_records", API_KEY)

	resBytes := []byte(res)
	var jsonRes ZoneRecordsResult
	err := json.Unmarshal(resBytes, &jsonRes)

	if err != nil {
		return []Record{}, err
	}

	return jsonRes.Result, nil
}

func getCurrentIpAddress() (string, error) {
	response := makeRequest(IP_API_ENDPOINT, "")

	parts := strings.Split(response, "ip=")
	ip := strings.Split(parts[1], "\n")[0]

	return ip, nil
}

func makeRequest(URL string, API_KEY string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer "+API_KEY)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Err is", err)
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	response := string(resBody)

	return response
}

func makePatch(URL string, API_KEY string, data []byte) string {
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", URL, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+API_KEY)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Err is", err)
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	response := string(resBody)

	return response
}

func printConig(config *config.Config) {
	log.Print("Config:")
	log.Print("- Zones - ")
	log.Print(config.Zones)
	log.Print("- DDNS Comment - ")
	log.Print(config.DDNSComment)
}

func run(config *config.Config) {
	currentIpAddress, err := getCurrentIpAddress()

	log.Print("Current IP: " + currentIpAddress)

	if err != nil {
		log.Panic("Couldn't get current IP address - exiting")
	}

	for _, zone := range config.Zones {
		records, err := getRecordsByZoneId(zone, config.ApiKey)

		if err != nil {
			log.Print("Couldn't get records for zone: " + zone)
			continue
		}

		for _, record := range records {
			if record.Comment != config.DDNSComment || record.Content == currentIpAddress {
				continue
			}

			err = updateRecord(zone, record.Id, currentIpAddress, config.ApiKey)
			if err != nil {
				log.Print("Couldnt update Record " + record.Name)
				continue
			}

			log.Print("Updated Record " + record.Name)
		}
	}
}
