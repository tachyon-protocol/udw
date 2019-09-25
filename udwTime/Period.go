package udwTime

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"
)

type PeriodHour struct {
	Start int
	End   int
}

type Period struct {
	Start time.Time
	End   time.Time
}

type DateStringPeriod struct {
	StartInclude string
	EndInclude   string
}

type PeriodGetter interface {
	GetPeriod() Period
}

func getReflectTypePeriodGetter() reflect.Type {
	return reflect.TypeOf((*PeriodGetter)(nil)).Elem()
}

func (p Period) IsIn(t time.Time) bool {
	if t.After(p.End) {
		return false
	}
	if t.Before(p.Start) {
		return false
	}
	return true
}

func (p Period) IsZero() bool {
	return p.Start.IsZero() && p.End.IsZero()
}

func (p Period) HasData() bool {
	return p.IsZero() == false
}

func (p Period) IsValid() bool {
	if p.End.Before(p.Start) || p.IsZero() == true {
		return false
	} else {
		return true
	}
}

func leq(a, b time.Time) bool {
	return a.Before(b) || a.Equal(b)
}

func (p Period) Overlaps(b Period) bool {
	return (leq(b.Start, p.Start) && leq(p.Start, b.End)) ||
		(leq(p.Start, b.Start) && leq(b.Start, p.End))
}

func NewPeriod(Start time.Time, End time.Time) (period Period, err error) {
	if Start.After(End) {
		err = fmt.Errorf("[udwTime.NewPeriod] Start.After(End) Start:%s End:%s", Start, End)
		return
	}
	return Period{Start: Start, End: End}, nil
}

func NewPeriodPassed24Hours() Period {
	now := time.Now()
	return Period{
		Start: now.Add(-24 * time.Hour),
		End:   now,
	}
}

func MustNewPeriod(Start time.Time, End time.Time) (period Period) {
	period, err := NewPeriod(Start, End)
	if err != nil {
		panic(err)
	}
	return period
}

func GetPeriodFromSortedList(t time.Time, SortedList []Period) (index int, ok bool) {
	n := len(SortedList)
	i := sort.Search(n, func(i int) bool {
		return SortedList[i].End.After(t)
	})
	if i == n {
		return 0, false
	}
	if !(SortedList[i].Start.Before(t) || SortedList[i].Start.Equal(t)) {
		return 0, false
	}
	return i, true
}

func GetPeriodFromGenericSortedList(t time.Time, SortedList interface{}) (index int, err error) {
	reflectList := reflect.Indirect(reflect.ValueOf(SortedList))
	if reflectList.Kind() != reflect.Slice && reflectList.Kind() != reflect.Array {
		panic(fmt.Errorf("[GetPeriodFromGenericSortedList] need array or slice get %s", reflectList.Kind().String()))
	}
	if !reflectList.Type().Elem().Implements(getReflectTypePeriodGetter()) {
		panic(fmt.Errorf("[GetPeriodFromGenericSortedList] need elem implement 'PeriodGetter' get %s",
			reflectList.Elem().Type().Name()))
	}
	n := reflectList.Len()
	i := sort.Search(n, func(i int) bool {
		return reflectList.Index(i).Interface().(PeriodGetter).GetPeriod().End.After(t)
	})
	if i == n {
		return 0, errors.New("not found")
	}
	if !reflectList.Index(i).Interface().(PeriodGetter).GetPeriod().Start.Before(t) {
		return 0, errors.New("not found")
	}
	return i, nil
}

type PeriodSlice []Period

func (p PeriodSlice) Len() int {
	return len(p)
}

func (p PeriodSlice) Less(i, j int) bool {
	return p[i].Start.Before(p[j].Start)
}
func (p PeriodSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func PeriodSort(p []Period) {
	sort.Sort(PeriodSlice(p))
}

type PeriodListInterface interface {
	Len() int
	GetPeriodAtIndex(i int) Period
	Swap(i, j int)
}

func SelectPeriodFromSortedPeriodList(t time.Time, PeriodList PeriodListInterface) (index int, ok bool) {
	n := PeriodList.Len()
	i := sort.Search(n, func(i int) bool {
		return PeriodList.GetPeriodAtIndex(i).End.After(t)
	})
	if i == n {
		return 0, false
	}
	currentPeriod := PeriodList.GetPeriodAtIndex(i)
	if !(currentPeriod.Start.Before(t) || currentPeriod.Start.Equal(t)) {
		return 0, false
	}
	return i, true
}

type periodListSorter struct {
	PeriodListInterface
}

func (p periodListSorter) Less(i, j int) bool {
	if p.GetPeriodAtIndex(i).End.Equal(p.GetPeriodAtIndex(j).End) {
		return p.GetPeriodAtIndex(i).Start.Before(p.GetPeriodAtIndex(j).Start)
	}
	return p.GetPeriodAtIndex(i).End.Before(p.GetPeriodAtIndex(j).End)
}
func PeriodListSort(PeriodList PeriodListInterface) {
	sort.Sort(periodListSorter{PeriodList})
}

type PeriodList []PeriodListElement

type PeriodListElement struct {
	Period      Period
	OriginIndex int
}

func (p PeriodList) Len() int {
	return len(p)
}
func (p PeriodList) Less(i, j int) bool {
	if p[i].Period.End.Equal(p[j].Period.End) {
		return p[i].Period.Start.Before(p[j].Period.Start)
	}
	return p[i].Period.End.Before(p[j].Period.End)
}
func (p PeriodList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PeriodList) Sort() {
	sort.Sort(p)
}

func (p PeriodList) SelectFromTime(t time.Time) (OriginIndex int, ok bool) {
	n := p.Len()
	i := sort.Search(n, func(i int) bool {
		return p[i].Period.End.After(t)
	})
	if i == n {
		return 0, false
	}
	currentPeriod := p[i]
	if !(currentPeriod.Period.Start.Before(t) || currentPeriod.Period.Start.Equal(t)) {
		return 0, false
	}
	return p[i].OriginIndex, true
}

func (p PeriodList) IsIncludeOverlap() bool {
	for i := 0; i < len(p)-1; i++ {
		if p[i].Period.Overlaps(p[i+1].Period) {
			return true
		}
	}
	return false
}
func (p PeriodList) SetAtIndex(period Period, i int) {
	p[i] = PeriodListElement{Period: period, OriginIndex: i}
}

func (p PeriodList) Append(period Period) PeriodList {
	return append(p, PeriodListElement{Period: period, OriginIndex: len(p)})
}
