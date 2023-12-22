package domain

import (
	"fmt"
	"time"
)

type ArchiverStatus int

const (
	Waiting ArchiverStatus = iota
	Running
	Complete
)

func (a ArchiverStatus) String() string {
	switch a {
	case 0:
		return "Waitng"
	case 1:
		return "Running"

	case 2:
		return "Complete"

	default:
		return "Unkown"
	}
}

// Archiver represents the contact archive functionality.
type Archiver struct {
	status   ArchiverStatus
	progress float64
}

func NewArchiver() *Archiver {
	return &Archiver{status: Waiting}
}

func (a *Archiver) Run() error {

	if a.status != Waiting {
		return fmt.Errorf("archive job already in progress or completed")
	}

	a.status = Running
	// Simulate some progress using goroutine
	go func() {
		defer func() {
			a.status = Complete
			a.progress = 1.0
		}()

		for i := 0; i <= 100; i += 10 {
			time.Sleep(500 * time.Millisecond)
			a.progress = float64(i) / 100.0
			fmt.Printf("Progress: %.0f%%\n", a.progress*100)
		}
	}()

	return nil
}
func (a *Archiver) Status() ArchiverStatus {

	return a.status
}

func (a *Archiver) Reset() {

	a.status = Waiting
	a.progress = 0
	fmt.Println("Archive job canceled.")
}

func (a *Archiver) Progress() float64 {

	return a.progress
}
