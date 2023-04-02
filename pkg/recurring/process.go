package recurring

import (
	"fmt"
	"time"
)

type (
	ProcessHandler   func()
	RecurringProcess struct {
		name      string
		interval  time.Duration
		startTime time.Time
		handler   func()
		stop      chan struct{}
	}
)

func NewRecurringProcess(
	name string,
	interval time.Duration,
	startTime time.Time,
	handler func(),
	stop chan struct{},
) *RecurringProcess {
	return &RecurringProcess{
		name:      name,
		interval:  interval,
		startTime: startTime,
		handler:   handler,
		stop:      stop,
	}
}

func (p *RecurringProcess) Name() string {
	return p.name
}

func (p *RecurringProcess) Stop() chan struct{} {
	return p.stop
}

func (p *RecurringProcess) Start() {
	startTicker := &time.Timer{}
	ticker := &time.Ticker{C: nil}
	defer func() { ticker.Stop() }()

	if p.startTime.Before(time.Now()) {
		p.startTime = time.Now()
	}
	startTicker = time.NewTimer(time.Until(p.startTime))

	for {
		select {
		case <-startTicker.C:
			ticker = time.NewTicker(p.interval)
			fmt.Println("Starting recurring process")
			p.handler()
		case <-ticker.C:
			fmt.Println("Next run")
			p.handler()
		case <-p.stop:
			fmt.Println("Stoping recurring process")
			return
		}
	}
}

func (p *RecurringProcess) Cancel() {
	close(p.stop)
}
