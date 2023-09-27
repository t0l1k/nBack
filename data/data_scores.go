package data

import (
	"container/list"
	"fmt"
	"sort"
	"time"

	"github.com/t0l1k/eui"
)

type ScoreData struct {
	Dt         time.Time
	Games, Max int
	Avg        float64
}

func (s *ScoreData) String() string {
	// dtFormat := "2006-01-02"
	layout := "02 Jan 2006"
	return fmt.Sprintf("%v %v:%v %v:%v %v:%v",
		s.Dt.Format(layout),
		eui.GetLocale().Get("wordGames"),
		s.Games,
		eui.GetLocale().Get("wordMax"),
		s.Max,
		eui.GetLocale().Get("wordAvg"),
		s.Avg)
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
		v := GetDb().ScoresData[k]
		idx.PushBack(i)
		maxs.PushBack(v.Max)
		averages.PushBack(v.Avg)
		if v.Games > 0 {
			strs.PushBack(v.String())
		} else {
			strs.PushBack("")
		}
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
