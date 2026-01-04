package worker

import (
	"log"
	"time"

	"AP_WebServer/internal/server"
)

func StartWorker(s *server.Server, stop <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("[STATS] requests=%d keys=%d",
				s.RequestCount(), s.KeyCount())
		case <-stop:
			log.Println("[WORKER] stopped")
			return
		}
	}
}
