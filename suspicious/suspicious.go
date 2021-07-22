package suspicious

import (
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
