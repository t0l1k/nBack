package data

import (
	"container/list"
	"fmt"
	"image/color"
	"math"
	"sort"
	"strconv"
	"time"

	ui "github.com/t0l1k/eui"
)

type TodayGamesData map[int]*GameData

func (t *TodayGamesData) GetToday() string {
	return time.Now().Format("2006.01.02")
}

func (t *TodayGamesData) GetCount() (count int) {
	return len(*t)
}

func (t *TodayGamesData) GetMax() (max int) {
	for _, v := range *t {
		if v.Level > max {
			max = v.Level
		}
	}
	return max
}

func (t *TodayGamesData) GetAvg() (sum float64) {
	for _, v := range *t {
		sum += float64(v.Level)
	}
	if t.GetCount() > 0 {
		sum /= float64(len(*t))
		return math.Round(sum*100) / 100
	}
	return 0
}

func (t *TodayGamesData) getGamesTimeDuraton() (result string) {
	if t.GetCount() == 0 {
		return
	}
	dtFormat := "2006-01-02 15:04:05.000"
	var durration time.Duration
	for _, v := range *t {
		dtBeg, err := time.Parse(dtFormat, v.DtBeg)
		if err != nil {
			panic(err)
		}
		dtEnd, err := time.Parse(dtFormat, v.DtEnd)
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
		v := GetDb().TodayData[k]
		gameNr.PushBack(k)
		level.PushBack(v.Level)
		result := float64(v.Percent)*0.01 + float64(v.Level)
		levelValue.PushBack(result)
		percents.PushBack(v.Percent)
		colors.PushBack(v.BgColor())
		moves := float64(v.Moves)
		totalmoves := float64(v.Totalmoves)
		percentMoves := moves * 100 / totalmoves
		lvlMoves := float64(v.Level) * percentMoves / 100
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
	advPercent := ui.GetPreferences().Get("threshold advance").(int)
	advCount := (*ui.GetPreferences()).Get("manual advance").(int)
	lastLvl := GetDb().TodayData[len(keys)].Level
	ok := false
	for i := len(keys); i > 0; i-- {
		v := GetDb().TodayData[i]
		if advCount == 0 || !v.Manual && count < advCount {
			return false, false, count
		} else if v.Manual && v.Percent >= advPercent && v.Level == lastLvl {
			count++
			ok = true
			if count == advCount {
				return true, ok, count
			}
		} else if v.Manual && v.Percent < advPercent && v.Level == lastLvl {
			ok = true
			return false, ok, count
		}
		if v.Level != lastLvl {
			return false, ok, count
		}
		lastLvl = v.Level
	}
	return count >= advCount, ok, count
}

func (t *TodayGamesData) ListShortStr() (strs []string, clrs []color.Color) {
	keys := make([]int, 0)
	for k := range *t {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, v := range keys {
		strs = append(strs, GetDb().TodayData[v].ShortStr())
		clrs = append(clrs, GetDb().TodayData[v].BgColor())
	}
	return
}

func (t *TodayGamesData) LongStr() (str string) {
	keys := make([]int, 0)
	for k := range *t {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for i, v := range keys {
		str += "#" + strconv.Itoa(i+1)
		str += GetDb().TodayData[v].String()
		if i < len(keys)-1 {
			str += "\n"
		}
	}
	return str
}
func (t *TodayGamesData) String() string {
	s := fmt.Sprintf("%v", t.GetToday())
	if t.GetCount() > 0 {
		s = fmt.Sprintf("%v #%v %v:%v, %v:%v [%v]",
			t.GetToday(),
			t.GetCount(),
			ui.GetLocale().Get("wordMax"),
			t.GetMax(),
			ui.GetLocale().Get("wordAvg"),
			t.GetAvg(),
			t.getGamesTimeDuraton(),
		)
	}
	return s
}
