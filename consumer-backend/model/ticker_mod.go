package model

import (
	"strings"
	"time"

	"github.com/thoas/go-funk"

	"github.com/agflow/tools/agstring"
	"github.com/agflow/tools/agtime"
)

const (
	SlowRefresh string = "3perweek"
	FastRefresh string = "Always"
)

type TickerMod struct {
	commodities     []string
	exportCountries []string
	incoterms       []string
	frequency       string
}

func (tck *TickerMod) ContainsCommodity(s string) bool {
	return funk.ContainsString(tck.commodities, "All") || funk.ContainsString(tck.commodities, s)
}

func (tck *TickerMod) ContainsIncoterm(s string) bool {
	return funk.ContainsString(tck.incoterms, "All") || funk.ContainsString(tck.incoterms, s)
}

func (tck *TickerMod) ContainsExportCountry(s string) bool {
	return funk.ContainsString(tck.exportCountries, "All") || funk.ContainsString(tck.exportCountries, s)
}

func (tck *TickerMod) ContainsAnyExportCountries(ss ...string) bool {
	return funk.ContainsString(tck.exportCountries, "All") || agstring.ContainsAny(tck.exportCountries, ss...)
}

func (tck *TickerMod) IsRefreshTimeValid(now time.Time, t agtime.NullableTime) bool {
	if tck.frequency == FastRefresh {
		return true
	}

	until := agtime.NewNullTime(now)
	if now.Weekday() == 2 || now.Weekday() == 4 {
		until = until.AddDate(0, 0, -1)
	}
	return until.AfterEqual(t)
}

func (tck *TickerMod) IsValidItem(now time.Time, i Item) bool {
	exportCountries := strings.Split(i.ExtExportCountry(), "/")
	isValidCommodity := tck.ContainsCommodity(i.ExtCommodity())
	isValidExportCountry := i.IsExportCountryValid() || tck.ContainsAnyExportCountries(exportCountries...)
	isValidIncoterm := i.IsIncotermValid() || tck.ContainsIncoterm(i.ExtIncoterm())
	isValidValidOn := i.IsValidOnValid() || tck.IsRefreshTimeValid(now, i.ExtValidOn())
	return isValidCommodity && isValidExportCountry && isValidIncoterm && isValidValidOn
}

// toTickerMod parses a rule string into ticker modifiers
// input example: "PDTrends-canSeeArgentina,Brazil/Corn/CPT/Always"
func toTickerMod(mod, s string) (*TickerMod, bool) {
	s = strings.ReplaceAll(s, mod, "")
	ruleDetails := strings.Split(s, "/")
	if len(ruleDetails) != 4 {
		return nil, false
	}

	tckMod := TickerMod{}

	// countries
	countries := strings.Split(ruleDetails[0], ",")
	tckMod.exportCountries = countries

	// commodity
	commodities := strings.Split(ruleDetails[1], ",")
	tckMod.commodities = commodities

	// incoterms
	incoterms := strings.Split(ruleDetails[2], ",")
	tckMod.incoterms = incoterms

	tckMod.frequency = ruleDetails[3]
	return &tckMod, true
}

func ToTickerMod(mod string, ls ...string) ([]*TickerMod, error) {
	tckMods := make([]*TickerMod, 0)
	for _, s := range ls {
		tckMod, ok := toTickerMod(mod, s)
		if !ok {
			continue
		}
		tckMods = append(tckMods, tckMod)
	}
	return tckMods, nil
}
