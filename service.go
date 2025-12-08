package main

import (
	"log"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

type SyncService struct{}

func (m *SyncService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	status <- svc.Status{State: svc.StartPending}
	go StartSyncLoop()
	status <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	for req := range r {
		switch req.Cmd {
		case svc.Interrogate:
			status <- req.CurrentStatus
		case svc.Stop, svc.Shutdown:
			status <- svc.Status{State: svc.StopPending}
			log.Println("Service arrêté.")
			return false, 0
		}
	}
	return false, 0
}

func runWindowsService() {
	err := svc.Run("SyncAppService", &SyncService{})
	if err != nil {
		log.Printf("Erreur service Windows : %v", err)
		debug.Run("SyncAppService", &SyncService{})
	}
}
