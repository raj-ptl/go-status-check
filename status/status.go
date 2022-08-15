package status

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type websiteStatus struct {
	URL         string
	LastChecked time.Time
}

type StatusChecker interface {
	Check(ctx context.Context, name string) (status bool, err error)
}

type HttpChecker struct {
}

func InitializeMap() {
	fmt.Println("InitializeMap() invoked ...")
	websiteMap := make(map[string]*websiteStatus)
	websiteMap["test.com"] = &websiteStatus{
		URL:         "test.com",
		LastChecked: time.Now(),
	}

	displayMap(&websiteMap)

}

func (h HttpChecker) Check(ctx context.Context, name string) (status bool, err error) {

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

	return true, nil
}

func displayMap(m *map[string]*websiteStatus) {
	for _, v := range *m {
		fmt.Printf("%+v\n", v)
	}
}
