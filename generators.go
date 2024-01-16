package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type NumberGenerator struct{
	counter int
	mutex sync.Mutex
}

func NewNumberGenerator() *NumberGenerator {
	return &NumberGenerator{}
}

func (g *NumberGenerator) Generate() string {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.counter++
	return fmt.Sprintf("%d", g.counter)
}


const (
	sequenceBits int64 = 22
	baseEpoch int64 = 1704067200000
	maxSequence = int64(-1) ^ (int64(-1) << sequenceBits)
	timeLeft int64 = 22
)
type SnowflakeIDGenerator struct {
	mutex sync.Mutex
	lastTimestamp int64
	sequence int64
}

func NewSnowflakeIDGenerator() *SnowflakeIDGenerator {
	return &SnowflakeIDGenerator{
		lastTimestamp: 0,
		sequence: 0,
	}
}

func (g *SnowflakeIDGenerator) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (g *SnowflakeIDGenerator) Generate() string {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	id, err := g.nextID()
	if err != nil {
		log.Fatal(err)
	}
	stringID := strconv.FormatUint(id, 36)

	return stringID
}

func (g *SnowflakeIDGenerator) nextID() (uint64, error) {
	timestamp := g.getMilliSeconds()
	if timestamp < g.lastTimestamp {
		return 0, fmt.Errorf("timestamp less than last timestamp")
	}

	if g.lastTimestamp == timestamp {
		g.sequence = (g.sequence + 1) & maxSequence

		if g.sequence == 0 {
			for timestamp <= g.lastTimestamp {
				timestamp = g.getMilliSeconds()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastTimestamp = timestamp

	id := ((timestamp - baseEpoch) << int64(timeLeft)) | g.sequence

	return uint64(id), nil
}