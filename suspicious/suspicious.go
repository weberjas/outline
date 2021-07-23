package suspicious

import (
	"fmt"
	"time"
)

type TimeStore struct {
	name      string
	beginning time.Time
	ending    time.Time
}

func (t TimeStore) ExtendLife(ext time.Duration) {
	t.ending.Add(ext)
}

func (t TimeStore) unexportedFunction() {
	fmt.Println("I'm just here for testing!")
}
