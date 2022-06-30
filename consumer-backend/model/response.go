package model

import (
	"errors"
	"reflect"
	"time"

	"github.com/agflow/tools/agtime"
)

type Item interface {
	ExtCommodity() string
	ExtExportCountry() string
	ExtIncoterm() string
	ExtValidOn() agtime.NullableTime
	IsCommodityValid() bool
	IsExportCountryValid() bool
	IsIncotermValid() bool
	IsValidOnValid() bool
	ToogleActive()
}

type Collection []Item

type Collections []Collection

type Response interface {
	GetCollections() Collections
	SetCollections(...Collection) interface{}
	Prefix() string
}

func (cll Collection) FilterTickerConstraints(now time.Time, tckMods ...*TickerMod) Collection {
	result := make(Collection, 0)
	for i := range cll {
		if IsValid(now, cll[i], tckMods...) {
			result = append(result, cll[i])
		}
	}
	return result
}

func (clls Collections) FilterTickerConstraints(now time.Time, tckMods ...*TickerMod) Collections {
	result := make(Collections, 0)
	for i := range clls {
		result = append(result, clls[i].FilterTickerConstraints(now, tckMods...))
	}
	return result
}

func IsValid(now time.Time, item Item, tckMods ...*TickerMod) bool {
	for i := range tckMods {
		if tckMods[i].IsValidItem(now, item) {
			return true
		}
	}
	return false
}

func ToCollection(cs interface{}) (Collection, error) {
	if reflect.TypeOf(cs).Kind() != reflect.Slice {
		return nil, errors.New("can't convert type to Item")
	}

	s := reflect.ValueOf(cs)
	result := make([]Item, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = s.Index(i).Interface().(Item)
	}
	return result, nil
}

func ToCollections(rawClls ...interface{}) (Collections, error) {
	result := make(Collections, len(rawClls))
	for i := range rawClls {
		cll, err := ToCollection(rawClls[i])
		if err != nil {
			return nil, err
		}
		result[i] = cll
	}
	return result, nil
}

func (cll Collection) MarkByTickerConstraints(now time.Time, tckMods ...*TickerMod) Collection {
	for i := range cll {
		if IsValid(now, cll[i], tckMods...) {
			cll[i].ToogleActive()
		}
	}
	return cll
}

func (clls Collections) MarkByTickerConstraints(now time.Time, tckMods ...*TickerMod) Collections {
	result := make(Collections, 0)
	for _, cll := range clls {
		result = append(result, cll.MarkByTickerConstraints(now, tckMods...))
	}
	return result
}
