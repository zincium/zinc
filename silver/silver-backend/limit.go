package main

import "sync"

var (
	DefaultAccessLimit = 20
)

type Limiter struct {
	maximumAccessLimit int
	table              map[int64]int
	mtx                sync.Mutex
}

func NewLimiter(maximumAccessLimit int) *Limiter {
	if maximumAccessLimit <= 0 {
		maximumAccessLimit = DefaultAccessLimit
	}
	return &Limiter{maximumAccessLimit: maximumAccessLimit, table: make(map[int64]int)}
}

func (lm *Limiter) Inc(id int64) bool {
	lm.mtx.Lock()
	defer lm.mtx.Unlock()
	v, ok := lm.table[id]
	if ok {
		if v > lm.maximumAccessLimit {
			return false
		}
		v++
	}
	lm.table[id] = v
	return true
}

func (lm *Limiter) Dec(id int64) int {
	lm.mtx.Lock()
	defer lm.mtx.Unlock()
	v, ok := lm.table[id]
	if !ok {
		return 0
	}
	v--
	if v <= 0 {
		delete(lm.table, id)
		return 0
	}
	lm.table[id] = v
	return v
}
