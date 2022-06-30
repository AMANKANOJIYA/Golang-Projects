package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func startHealthCheck() {
	s := gocron.NewScheduler(time.Local)
	for _, host := range serverList {
		_, err := s.Every(2).Seconds().Do(func(s *server) {
			healthy := s.checkHealth()
			if healthy {
				log.Printf("%s is Healthy", s.Name)
			} else {
				log.Printf("%s is Not Healthy", s.Name)
			}
		}, host)
		if err != nil {
			log.Fatalln(err)
		}
	}
	s.StartAsync()
}
