package status

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/raj-ptl/go-status-check/models"
)

type StatusChecker interface {
	Check(ctx context.Context, name string) (status bool, err error)
}

type HttpChecker struct {
}

var hc HttpChecker

var WebsiteMap = make(map[string]*models.WebsiteStatus)
var WebsiteMapMutex = sync.RWMutex{}

func ExposeMap() *map[string]*models.WebsiteStatus {
	return &WebsiteMap
}

func (h HttpChecker) Check(ctx context.Context, name string) (status string, err error) {

	if strings.Split(name, ":")[0] != "http" || strings.Split(name, ":")[0] != "https" {
		name = "http://" + name
	}

	res, err := http.Get(name)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != 200 {
		return "DOWN", nil
	}

	return "UP", nil
}

func UpdateSingleSite(url string, ch chan int) {

	logFile, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("ERROR while opening log file : %v", err)
	}

	log.SetOutput(logFile)

	defer logFile.Close()

	status, _ := hc.Check(context.TODO(), url)
	WebsiteMapMutex.Lock()
	WebsiteMap[url] = &models.WebsiteStatus{
		URL:         url,
		Status:      status,
		LastChecked: time.Now(),
	}
	WebsiteMapMutex.Unlock()

	log.Printf("INFO : Updated status for site : %s -> %s", url, status)

	<-ch
}

func UpdateSingleSiteSynchronous(url string) {
	status, _ := hc.Check(context.TODO(), url)
	WebsiteMapMutex.Lock()
	WebsiteMap[url] = &models.WebsiteStatus{
		URL:         url,
		Status:      status,
		LastChecked: time.Now(),
	}
	WebsiteMapMutex.Unlock()
}

func UpdateAllSites() {
	ch := make(chan int)

	for key := range WebsiteMap {
		go UpdateSingleSite(key, ch)
	}

}

/*
pollingRate -> Seconds
*/
func PollUpdateAllSites(pollingRate float32) {
	for {
		UpdateAllSites()
		fmt.Printf("Sleeping for %v seconds\n", pollingRate)
		time.Sleep(time.Duration(pollingRate) * time.Second)
	}
}
