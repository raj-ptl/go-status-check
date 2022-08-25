package status

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	status, _ := hc.Check(context.TODO(), url)
	WebsiteMapMutex.Lock()
	WebsiteMap[url] = &models.WebsiteStatus{
		URL:         url,
		Status:      status,
		LastChecked: time.Now(),
	}
	WebsiteMapMutex.Unlock()
	<-ch
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
		fmt.Printf("Invoking update all sites\n")
		UpdateAllSites()
		fmt.Printf("Sleeping for %v seconds\n", pollingRate)
		time.Sleep(time.Duration(pollingRate) * time.Second)
	}
}
