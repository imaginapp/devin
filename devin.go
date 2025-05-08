package devin

import (
	"time"

	"golang.org/x/exp/rand"
)

const OpenInviteCodeGroupID = "open"
const OpenInviteMinimumBalance = 30

func RandomSleep(minMs, maxMs int) {
	sleepMs := rand.Intn(maxMs-minMs) + minMs
	time.Sleep(time.Duration(sleepMs) * time.Millisecond)
}

func RandomDelay() {
	RandomSleep(1000, 5000)
}
