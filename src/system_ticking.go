package main

import (
	"sync"
	"time"

	"github.com/TunnelWork/Ulysses/src/internal/logger"
)

type tickEventSignature uint64
type tickEvent func()

const (
	// Don't set high ticking frequency for a high (detailed) logging level.
	tickEventPeriodMsDefault uint16 = 1000 // 1s/tick.
)

var (
	tickEventMutex    *sync.Mutex // SGL here. Reason: Only 1 read goroutine as tickWorker. All other calls are write.
	tickEventMap      map[tickEventSignature]tickEvent
	tickEventPeriodMs uint16
	tickTicker        *time.Ticker
)

func registerTickEvent(tickEventSign tickEventSignature, tickEvt func()) {
	tickEventMutex.Lock()
	defer tickEventMutex.Unlock()

	if _, ok := tickEventMap[tickEventSign]; ok {
		logger.Debug("registerTickEvent(): repeated signature ", tickEventSign)
	}
	tickEventMap[tickEventSign] = tickEvt
}

func removeTickEvent(tickEventSign tickEventSignature) {
	tickEventMutex.Lock()
	defer tickEventMutex.Unlock()

	if _, ok := tickEventMap[tickEventSign]; ok {
		logger.Debug("removeTickEvent(): removing event with signature ", tickEventSign)
		delete(tickEventMap, tickEventSign)
	} else {
		logger.Debug("removeTickEvent(): signature not found ", tickEventSign)
	}
}

func initSystemTicking() {
	tickEventMutex = &sync.Mutex{}
	if masterConfig.Sys.SystemTickPeriodMillisecond == 0 {
		tickEventPeriodMs = tickEventPeriodMsDefault
	} else {
		tickEventPeriodMs = masterConfig.Sys.SystemTickPeriodMillisecond
	}
	tickEventMap = map[tickEventSignature]tickEvent{}
}

func startSystemTicking() {
	if tickEventPeriodMs == 0 {
		logger.Fatal("startSystemTicking(): ticker uninitialized. call initSystemTicking().")
	}

	tickTicker = time.NewTicker(time.Duration(tickEventPeriodMs) * time.Millisecond)
	globalWaitGroup.Add(1)

	go func() {
		defer globalWaitGroup.Done()
		for {
			select {
			case <-globalExitChannel:
				return
			case <-tickTicker.C:
				singleTick()
			}
		}
	}()
}

func singleTick() {
	tickEventMutex.Lock()
	defer tickEventMutex.Unlock()
	logger.Debug("singleTick(): now ticking")
	for signature, evt := range tickEventMap {
		logger.Debug("singleTick(): execute event with signature ", signature)
		evt()
	}
}
