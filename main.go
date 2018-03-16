package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
)

func worker(id int, jobs <-chan *http.Request, results chan<- *http.Response) {
    for req := range jobs {
        fmt.Printf("\nworker %v: started %v", id, req)
        res, err := http.DefaultClient.Do(req)
        fmt.Printf("\nworker %v: finished", id, req)
        if err != nil {
            fmt.Printf("error processing request: %+v\n", err)
        } else {
            results <- res
        }
    }
}

func main() {
    workerFlag := flag.Int("worker-count", 300, "number of workers")
    requestCountFlag := flag.Int("request-count", 10000, "number of requests")
    urlFlag := flag.String("request-url", "https://google.com", "url to request")
    flag.Parse()

    workChan := make(chan *http.Request, *requestCountFlag)
    resultChan := make(chan *http.Response)

    go func() {
        defer close(workChan)
        for i := 0; i < *requestCountFlag; i++ {
            req, err := http.NewRequest("GET", *urlFlag, nil)
            if err != nil {
                fmt.Printf("error creating request: %v\n", err)
            } else {
                workChan <- req
            }
        }
    }()
    for i := 0; i < *workerFlag; i++ {
        go worker(i, workChan, resultChan)
    }

    for i := 0; i < *requestCountFlag; i++ {
        res := <-resultChan
        body, _ := ioutil.ReadAll(res.Body)
        fmt.Printf("\nresult body: %+v", string(body))
        fmt.Printf("\nresult status: %+v", res.StatusCode)
    }
    close(resultChan)

}
