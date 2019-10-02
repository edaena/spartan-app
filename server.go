package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var appVersion = "0.1" //Default/fallback version
var instanceNum int

func getFrontpage(w http.ResponseWriter, r *http.Request) {
	t := time.Now().UTC()
	fmt.Fprintf(w, "Hi folks! I'm instance %d running version %s of your application at %s\n", instanceNum, appVersion, t.Format("2006-01-02 15:04:05"))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	recordMetrics()
	fmt.Fprintf(w, "%s\n", appVersion)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	instanceNum = rand.Intn(1000)
	recordError()
	http.HandleFunc("/", getFrontpage)
	http.HandleFunc("/health", health)
	http.HandleFunc("/version", getVersion)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//Sum : Simple summing function.
func Sum(x int, y int) int {
	return x + y
}

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

//Quick and dirty simulation error state
func recordError() {
	go func() {
		var divisor = 17 //Higher odd number means potential longer time to trigger error
		var errorFound = false
		var errorIterations = 5
		for {
			var randomNumberDividend = rand.Intn(1000)
			if errorFound && errorIterations > 0 {
				errorIterations--
				opsErrorState.Set(1)
				fmt.Printf("Set ErrorState to %v\n", 1)
			} else if errorFound && errorIterations == 0 {
				errorIterations = 5
				errorFound = false
				opsErrorState.Set(0)
				fmt.Printf("Set ErrorState to %v\n", 0)
			} else if randomNumberDividend%divisor == 0 {
				opsErrorState.Set(1)
				errorFound = true
				fmt.Printf("Set ErrorState to %v\n", 1)
			} else {
				opsErrorState.Set(0)
				fmt.Printf("Set ErrorState to %v\n", 0)
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "spartan_app_processed_ops_total",
		Help: "The total number of processed events",
	})
	//A simple numeric value that can go up and down and exposes a current state of the application.
	opsErrorState = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "spartan_app_error_state",
		Help: "Counter that counts the number of successes",
	})
)
