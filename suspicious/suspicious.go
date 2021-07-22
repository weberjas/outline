package suspicious

import (
	"time"
)

type TimeStore struct {
	name     string
	begining time.Time
	ending   time.Time
}
