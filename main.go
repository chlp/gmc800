package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/tarm/serial"
)

var (
	currentCPM uint32
	mu         sync.RWMutex
)

func main() {
	var fixedPort string
	if len(os.Args) >= 2 {
		fixedPort = os.Args[1]
	}

	go func() {
		for {
			var port string
			if fixedPort != "" {
				port = fixedPort
			} else {
				p, err := findSerialPort()
				if err != nil {
					log.Println("Serial Port not found:", err)
					time.Sleep(3 * time.Second)
					continue
				}
				port = p
				log.Println("Serial Port found:", port)
			}

			connectAndPollLoop(port)

			if fixedPort != "" {
				time.Sleep(3 * time.Second)
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		temps := getTemperatures()
		resp := map[string]interface{}{
			"cpm":  currentCPM,
			"temp": temps,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("HTTP-started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connectAndPollLoop(portName string) {
	log.Println("Connecting with Serial Port:", portName)
	s, err := serial.OpenPort(&serial.Config{
		Name:        portName,
		Baud:        115200,
		ReadTimeout: time.Second,
	})
	if err != nil {
		log.Println("Could not open port:", err)
		return
	}
	defer s.Close()

	log.Println("Success:", portName)

	for {
		err := requestAndUpdateCPM(s)
		if err != nil {
			log.Println("Wrong request CPM:", err)
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func requestAndUpdateCPM(s *serial.Port) error {
	_, err := s.Write([]byte("<GETCPM>>"))
	if err != nil {
		return fmt.Errorf("could not write: %w", err)
	}

	buf := make([]byte, 4)
	n, err := s.Read(buf)
	if err != nil {
		return fmt.Errorf("could not read: %w", err)
	}
	if n != 4 {
		return fmt.Errorf("wrong response: %d байт", n)
	}

	cpm := binary.BigEndian.Uint32(buf)

	mu.Lock()
	currentCPM = cpm
	mu.Unlock()

	return nil
}

func getTemperatures() map[string]float64 {
	result := make(map[string]float64)

	temps, err := host.SensorsTemperatures()
	if err != nil {
		result["error"] = -1
		return result
	}

	for _, t := range temps {
		if t.Temperature > 0 && t.Temperature < 100 {
			result[t.SensorKey] = t.Temperature
		}
	}
	return result
}

func findSerialPort() (string, error) {
	matches, err := filepath.Glob("/dev/tty.usbserial-*")
	if err != nil {
		return "", fmt.Errorf("could not find port: %w", err)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("could not find matching port")
	}
	return matches[len(matches)-1], nil
}
