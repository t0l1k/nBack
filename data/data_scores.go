package data

import (
	"container/list"
	"fmt"
	"sort"
	"time"

	"github.com/t0l1k/nBack/ui"
)

type ScoreData struct {
	Dt         time.Time
	Games, Max int
	Avg        float64
}

func (s *ScoreData) String() string {
	dtFormat := "2006-01-02"
	return fmt.Sprintf("%v %v:%v %v:%v %v:%v",
		s.Dt.Format(dtFormat),
		ui.GetLocale().Get("wordGames"),
		s.Games,
		ui.GetLocale().Get("wordMax"),
		s.Max,
		ui.GetLocale().Get("wordAvg"),
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
