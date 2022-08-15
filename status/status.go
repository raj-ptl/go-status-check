package status

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type WebsiteStatus struct {
	URL         string
	Status      string
	LastChecked time.Time
}

type StatusChecker interface {
	Check(ctx context.Context, name string) (status bool, err error)
}

type HttpChecker struct {
}

var hc HttpChecker

var WebsiteMap = make(map[string]*WebsiteStatus)

func ExposeMap() *map[string]*WebsiteStatus {
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

	fmt.Printf("%s:\n", name)
	fmt.Printf("Status : %s:\n", res.Status)
	fmt.Printf("Status Code : %d:\n", res.StatusCode)

	if res.StatusCode != 200 {
		return "DOWN", nil
	}

	return "UP", nil
}

func DisplayMap(m *map[string]*WebsiteStatus) {
	for _, v := range *m {
		fmt.Printf("%+v\n", v)
	}
}

func UpdateSingleSite(url string) {
	status, _ := hc.Check(context.TODO(), url)
	WebsiteMap[url] = &WebsiteStatus{
		URL:         url,
		Status:      status,
		LastChecked: time.Now(),
	}
}
