/**
 * 
 * aerTest.go
 * A simple multi-thread stress test for AER Service using connectorlog rethinkdb database
 *
 * Tomeu Capó 2017 (C)
 */

package main

import (
	"aerTest/domain"
	"aerTest/utils"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {

	numThreadsPtr := flag.Int("n", 4, "int")

	flag.Parse()

	if *numThreadsPtr <= 0 {
		log.Fatalln("n greater than zero expected")
	}

	conf := utils.ReadConfig("config.toml")
	confDb := utils.GetRDBConnectArguments(&conf)

	var channelRequests chan domain.Request = make(chan domain.Request)
	for _, requestType := range conf.AERConfig.RequestTypes {
		go utils.LogsReader(fmt.Sprintf("motor.%s", requestType), confDb, channelRequests)
	}

	fmt.Printf("Starting stress-test with %d threads ...\n", *numThreadsPtr)
	fmt.Printf("Log requests above %d ms ...\n", conf.AERConfig.RsPrintAbove)

	for t := 0; t < *numThreadsPtr; t++ {
		go func(c chan domain.Request) {
			var netTransport = &http.Transport{
				Dial: (&net.Dialer{
					Timeout: conf.AERConfig.Timeout * time.Second,
				}).Dial,
			}

			var client = &http.Client{
				Timeout:   conf.AERConfig.Timeout * time.Second,
				Transport: netTransport}

			urlService := conf.AERConfig.Endpoint
			buffRq := new(bytes.Buffer)

			var headers = http.Header{}
			headers.Set("Accept", "application/json")
			headers.Set("Content-Type", "application/json")

			for {
				rq := <-c

				buffRq.Reset()
				json.NewEncoder(buffRq).Encode(rq)
				rqStr := buffRq.String()

				req, _ := http.NewRequest("POST", urlService, buffRq)
				req.Header = headers

				tIni := time.Now().UnixNano() / int64(time.Millisecond)
				client.Do(req)
				tFi := time.Now().UnixNano() / int64(time.Millisecond)

				if tFi-tIni > conf.AERConfig.RsPrintAbove {
					log.Printf("[TX RQ] %s Total time = %d ms\n", rqStr, tFi-tIni)
				}
			}
		}(channelRequests)
	}

	var input string
	fmt.Scanln(&input)
}
