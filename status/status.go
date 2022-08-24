package status

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
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

func ExposeMap() *map[string]*models.WebsiteStatus {
	return &WebsiteMap
}

func InitializeMap() {
	fmt.Println("InitializeMap() invoked ...")

	DisplayMap(&WebsiteMap)

}

func (h HttpChecker) Check(ctx context.Context, name string) (status string, err error) {

	fmt.Println("Checking for : ", name)

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

func DisplayMap(m *map[string]*models.WebsiteStatus) {
	for _, v := range *m {
		fmt.Printf("%+v\n", v)
	}
}

func UpdateSingleSite(url string, ch chan int) {
	status, _ := hc.Check(context.TODO(), url)
	WebsiteMap[url] = &models.WebsiteStatus{
		URL:         url,
		Status:      status,
		LastChecked: time.Now(),
	}
	<-ch
}

/*
pollingRate -> Seconds
*/
func UpdateAllSites() {
	ch := make(chan int)

	for key := range WebsiteMap {
		go UpdateSingleSite(key, ch)
	}

}

func PollUpdateAllSites(pollingRate float32) {
	for {
		fmt.Printf("Invoking update all sites\n")
		UpdateAllSites()
		fmt.Printf("Sleeping for 10s\n")
		time.Sleep(time.Duration(pollingRate) * time.Second)
	}
}
