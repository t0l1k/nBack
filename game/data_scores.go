package game

import (
	"container/list"
	"fmt"
	"sort"
	"time"
)

type ScoreData struct {
	dt         time.Time
	games, max int
	avg        float64
}

func (s *ScoreData) String() string {
	dtFormat := "2006-01-02"
	return fmt.Sprintf("%v Игр:%v максимально:%v среднее:%v", s.dt.Format(dtFormat), s.games, s.max, s.avg)
}

type ScoresData map[int]*ScoreData

func (s *ScoresData) PlotData() (idx, maxs, averages, strs list.List) {

	keys := make([]int, 0)
	for k := range *s {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	i := 1
	for _, k := range keys {
		v := getDb().scoresData[k]
		idx.PushBack(i)
		maxs.PushBack(v.max)
		averages.PushBack(v.avg)
		strs.PushBack(v.String())
		i++
	}
	return
}

func (s *ScoresData) String() string {
	ss := ""
	for k, v := range *s {
		ss += fmt.Sprintf("%v [%v]\n", k, v)
	}
	return ss
}
