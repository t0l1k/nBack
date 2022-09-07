package game

import (
	"container/list"
	"fmt"
	"image/color"
	"math"
	"sort"
	"time"

	"github.com/t0l1k/nBack/ui"
)

type TodayGamesData map[int]*GameData

func (t *TodayGamesData) getToday() string {
	return time.Now().Format("2006.01.02")
}

func (t *TodayGamesData) getCount() (count int) {
	return len(*t)
}

func (t *TodayGamesData) getMax() (max int) {
	for _, v := range *t {
		if v.level > max {
			max = v.level
		}
	}
	return max
}

func (t *TodayGamesData) getAvg() (sum float64) {
	for _, v := range *t {
		sum += float64(v.level)
	}
	if t.getCount() > 0 {
		sum /= float64(len(*t))
		return math.Round(sum*100) / 100
	}
	return 0
}

func (t *TodayGamesData) getGamesTimeDuraton() (result string) {
	if t.getCount() == 0 {
		return
	}
	dtFormat := "2006-01-02 15:04:05.000"
	var durration time.Duration
	for _, v := range *t {
		dtBeg, err := time.Parse(dtFormat, v.dtBeg)
		if err != nil {
			panic(err)
		}
		dtEnd, err := time.Parse(dtFormat, v.dtEnd)
		if err != nil {
			panic(err)
		}
		durration += dtEnd.Sub(dtBeg)
	}
	d := durration.Round(time.Millisecond)
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	sec := d / time.Second
	d -= sec * time.Second
	ms := d / time.Millisecond
	if hours > 0 {
		result = fmt.Sprintf("%02v:%02v:%02v.%03v", int(hours), int(minutes), int(sec), int(ms))
	} else {
		result = fmt.Sprintf("%02v:%02v.%03v", int(minutes), int(sec), int(ms))
	}
	return
}

func (t *TodayGamesData) PlotTodayData() (gameNr, level, levelValue, percents, movesPerceent, colors list.List) {
	keys := make([]int, 0)
	for k := range *t {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		v := getDb().todayData[k]
		gameNr.PushBack(k)
		level.PushBack(v.level)
		result := float64(v.percent)*0.01 + float64(v.level)
		levelValue.PushBack(result)
		percents.PushBack(v.percent)
		colors.PushBack(v.BgColor())
		moves := float64(v.moves)
		totalmoves := float64(v.totalmoves)
		percentMoves := moves * 100 / totalmoves
		lvlMoves := float64(v.level) * percentMoves / 100
		movesPerceent.PushBack(lvlMoves)
	}
	return
}

func (t *TodayGamesData) getWinCountInManual() (bool, bool, int) {
	keys := make([]int, 0)
	for k := range *t {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	count := 0
	adv := (*ui.GetPreferences()).Get("manual advance").(int)
	lastLvl := getDb().todayData[len(keys)].level
	ok := false
	for i := len(keys); i > 0; i-- {
		v := getDb().todayData[i]
		if adv == 0 || !v.manual && count < adv {
			return false, false, count
		} else if v.manual && v.percent == 100 && v.level == lastLvl {
			count++
			ok = true
			if count == adv {
				return true, ok, count
			}
		} else if v.manual && v.percent < 100 && v.level == lastLvl {
			ok = true
			return false, ok, count
		}
		if v.level != lastLvl {
			return false, ok, count
		}
		lastLvl = v.level
	}
	return count >= adv, ok, count
}

func (t *TodayGamesData) ListShortStr() (strs []string, clrs []color.Color) {
	keys := make([]int, 0)
	for k := range *t {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, v := range keys {
		strs = append(strs, getDb().todayData[v].ShortStr())
		clrs = append(clrs, getDb().todayData[v].BgColor())
	}
	return
}

func (t *TodayGamesData) String() string {
	s := fmt.Sprintf("%v", t.getToday())
	if t.getCount() > 0 {
		s = fmt.Sprintf("%v #%v %v:%v, %v:%v [%v]",
			t.getToday(),
			t.getCount(),
			ui.GetLocale().Get("wordMax"),
			t.getMax(),
			ui.GetLocale().Get("wordAvg"),
			t.getAvg(),
			t.getGamesTimeDuraton(),
		)
	}
	return s
}
