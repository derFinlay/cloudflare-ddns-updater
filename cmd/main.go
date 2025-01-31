package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/derfinlay/ddns/config"
)

type Record struct {
	Id      string `json:"id"`
	Comment string `json:"comment"`
	Content string `json:"content"`
	Name    string `json:"name"`
	Zone    string `json:"zone_id"`
}

type ZoneRecordsResult struct {
	Result []Record `json:"result"`
}

var CLOUDFLARE_API_BASE_URL string = "https://api.cloudflare.com/client/v4/"
var IP_API_ENDPOINT string = "https://cloudflare.com/cdn-cgi/trace"
var NUM_ROUTINES int = 5

func main() {
	for {
		c, err := config.LoadConfig()

		if err != nil {
			log.Print("Error loading config - retrying in 10 seconds")
			time.Sleep(10 * time.Second)
			continue
		}

		current_IP, err := getCurrentIpAddress()
		if err != nil {
			log.Print("Couldn't get current IP address - retrying in 10 seconds")
			time.Sleep(10 * time.Second)
			continue
		}

		zoneChannel := make(chan string)
		go func() {
			for zone := range c.Zones {
				zoneChannel <- c.Zones[zone]
			}
			close(zoneChannel)
		}()

		log.Print("Starting Cloudflare Record updates...", time.Now())

		wg := sync.WaitGroup{}
		recordChannel := make(chan Record)
		for i := 0; i < NUM_ROUTINES; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for zone := range zoneChannel {
					records, err := getRecordsByZoneId(zone, c.ApiKey)
					if err != nil {
						log.Print("Couldn't get records for zone: " + zone)
						continue
					}
					for _, record := range records {
						recordChannel <- record
					}
				}
			}()
		}

		go func() {
			wg.Wait()
			close(recordChannel)
		}()

		var wg2 sync.WaitGroup
		for i := 0; i < NUM_ROUTINES; i++ {
			wg2.Add(1)
			go func() {
				defer wg2.Done()
				for record := range recordChannel {
					updateRecord(record.Zone, record, c.DDNSComment, current_IP, c.ApiKey)
				}
			}()
		}

		wg2.Wait()
		log.Print("All records updated")

		if c.UpdateInterval == 0 {
			os.Exit(0)
		}
		log.Print("-----------------------------")
		time.Sleep(time.Duration(c.UpdateInterval) * time.Second)
	}
}

func updateRecord(zone_id string, record Record, ddns_comment string, new_value string, API_KEY string) {
	if record.Comment != ddns_comment {
		log.Printf("Skipping record %s as it does not match the DDNS comment", record.Name)
		return
	}

	if record.Content == new_value {
		log.Printf("Skipping record %s as IP address is already up to date", record.Name)
		return
	}
	log.Print("Updating Record " + record.Id + " in Zone " + zone_id + " to " + new_value)
	body := []byte(`{"content": "` + new_value + `"}`)
	makeRequest(CLOUDFLARE_API_BASE_URL+"/zones/"+zone_id+"/dns_records/"+record.Id, http.MethodPatch, API_KEY, body)
	log.Print("Updated record " + record.Name + " to " + new_value)

}

func getRecordsByZoneId(zone_id string, API_KEY string) ([]Record, error) {
	log.Print("Getting records for zone: " + zone_id)
	res := makeRequest(CLOUDFLARE_API_BASE_URL+"/zones/"+zone_id+"/dns_records", http.MethodGet, API_KEY, nil)

	resBytes := []byte(res)
	var jsonRes ZoneRecordsResult
	err := json.Unmarshal(resBytes, &jsonRes)

	if err != nil {
		return []Record{}, err
	}

	return jsonRes.Result, nil
}

func getCurrentIpAddress() (string, error) {
	response := makeRequest(IP_API_ENDPOINT, http.MethodGet, "", nil)

	parts := strings.Split(response, "ip=")
	ip := strings.Split(parts[1], "\n")[0]
	log.Print("Current IP: " + ip)

	return ip, nil
}

func makeRequest(URL string, method string, API_KEY string, data []byte) string {
	client := &http.Client{}
	req, _ := http.NewRequest(method, URL, bytes.NewBuffer(data))
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
