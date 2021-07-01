package etcd

import (
	"sort"
	"time"
)

type Load interface {
	getAddr(data map[string]int64) string
}

var randomGlobal *Random

func GetRandom() *Random {
	if randomGlobal == nil {
		randomGlobal = &Random{}
	}
	return randomGlobal
}

type Random struct {
}

func (rl *Random) getAddr(data map[string]int64) string {
	addr := ""
	for k, v := range data {
		if v > time.Now().Unix() {
			addr = k
		} else {
			delete(data, k)
		}
	}
	return addr
}

var roundRobinGlobal *roundRobin

func GetRoundRobin() *roundRobin {
	if roundRobinGlobal == nil {
		roundRobinGlobal = &roundRobin{i: 0}
	}
	return roundRobinGlobal
}

type roundRobin struct {
	i int
}

func (rl *roundRobin) getAddr(data map[string]int64) string {
	var list []string
	for k, v := range data {
		if v > time.Now().Unix() {
			list = append(list, k)
		} else {
			delete(data, k)
		}
	}
	if list == nil {
		return ""
	} else {
		sort.Strings(list)
		if rl.i >= len(list) {
			rl.i = 1
		} else {
			rl.i++
		}
		return list[rl.i-1]
	}
}
