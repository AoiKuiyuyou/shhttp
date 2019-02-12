package main

import (
	"flag"
	"github.com/golang/glog"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func main() {
	hostname := flag.String("hostname", "0.0.0.0", "hostname to listen on")
	port := flag.Int("port", 2112, "port to listen on")
	cleanup := flag.Int("clean-interval", -1, "interval (hours) after which finished jobs are cleaned")
	location := flag.String("dir", "shhttp", "location to store the job data")
	revive := flag.Bool("revive", false, "Whether to revive previous running jobs if there are any")

	flag.Parse()

	jobStore := FileBasedJobStore{BasePath: path.Join(*location, "jobs")}
	jobStore.ensureDirectory()

	savedJobStore := FileBasedJobStore{BasePath: path.Join(*location, "saved")}
	savedJobStore.ensureDirectory()

	if *cleanup > 0 {
		interval := time.Duration(*cleanup) * time.Hour
		go func() {
			for {
				start := time.Now()
				jobStore.ClearFinished()
				elapsed := time.Since(start)
				if elapsed < interval {
					time.Sleep(interval - elapsed)
				}
			}
		}()
	}

	router := GetRouter(jobStore, savedJobStore, *revive)
	address := strings.Join([]string{*hostname, strconv.Itoa(*port)}, ":")
	glog.Infof("starting HTTP listener at %s", address)
	glog.Fatal(http.ListenAndServe(address, router))
}
