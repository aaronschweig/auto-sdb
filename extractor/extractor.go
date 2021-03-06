package extractor

import (
	"errors"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"

	"github.com/aaronschweig/auto-sdb/helpers"
)

type Extractor interface {
	WithContent(content string) Extractor
	ExtractBezeichnung() error
	ExtractSignalwort() error
	ExtractLagerklasse() error
	ExtractHPSaetze() error
	Extract() *SicherheitsdatenblattData
}

type SicherheitsdatenblattData struct {
	Bezeichnung string   `json:"bezeichnung"`
	Lagerklasse string   `json:"lagerklasse"`
	Signalwort  string   `json:"signalwort"`
	HSaetze     []string `json:"hSaezte"`
	PSaetze     []string `json:"pSaezte"`
}

var (
	Lagerklassen = [...]string{"1", "2A", "2B", "3", "4.1A", "4.1B", "4.2", "4.3", "5.1A", "5.1B", "5.1C", "5.2", "6.1A",
		"6.1B", "6.1C", "6.1D", "6.2", "7", "8A", "8B", "10", "11", "12", "13", "10-13"}
	hpSatzRegex       = regexp.MustCompile(`(?im)(\s?\+?\s?E?U?[HP][0-9]{3}[a-zA-Z]{0,2}){1,3}`)
	signalwortRegex   = regexp.MustCompile(`(?im)signalwort(.*)`)
	lagerklassenRegex = regexp.MustCompile(`(?im)lagerklasse(.*)`)
)

type DefaultExtractor struct {
	content string
	result  *SicherheitsdatenblattData
}

func NewDefaultExtractor() Extractor {
	return &DefaultExtractor{
		result: &SicherheitsdatenblattData{},
	}
}

func (e *DefaultExtractor) WithContent(content string) Extractor {
	e.content = content
	return e
}

func (e *DefaultExtractor) ExtractBezeichnung() error {
	bzRegex := regexp.MustCompile(`(?im)((handels?)?name|produktidentifikator)(\s*)\n(.*)`)

	matches := bzRegex.FindAllStringSubmatch(e.content, -1)

	if len(matches) == 0 {
		return errors.New("Could not extract bezeichnung.")
	}

	for i := range matches {
		for j := 1; j < len(matches[i]); j++ {
			candidate := matches[i][j]

			// ToLower
			candidate = strings.ToLower(candidate)

			// Strip Keyword
			candidate = strings.Replace(candidate, "produktidentifikator", "", -1)
			candidate = strings.Replace(candidate, "handelsname", "", -1)
			candidate = strings.Replace(candidate, ":", "", -1)

			// Trim
			candidate = strings.TrimSpace(candidate)

			// Überschriftszeile
			if len(candidate) == 0 {
				continue
			}
			e.result.Bezeichnung = candidate
			return nil
		}
	}

	return errors.New("Could not extract bezeichnung.")
}

func (e *DefaultExtractor) ExtractSignalwort() error {
	matches := signalwortRegex.FindAllString(e.content, -1)

	if len(matches) == 0 {
		return errors.New("Could not extract signalwort.")
	}

	for _, match := range matches {
		match = strings.ToLower(match)
		if strings.Contains(match, "gefahr") {
			e.result.Signalwort = "Gefahr"
			return nil
		}
		if strings.Contains(match, "achtung") {
			e.result.Signalwort = "Achtung"
			return nil
		}
	}
	return errors.New("Could not extract signalwort.")
}

func (e *DefaultExtractor) ExtractLagerklasse() error {
	matches := lagerklassenRegex.FindAllString(e.content, -1)

	if len(matches) == 0 {
		return errors.New("Could not extract lagerklasse.")
	}

	for _, match := range matches {
		if strings.Contains(match, "510") {
			match = strings.ReplaceAll(match, "510", "")
		}
		// Um Highest-Value Matches abzufangen wird von hinten durchiteriert, damit 1 nicht alle 1x Werte abfängt
		for i := len(Lagerklassen) - 1; i >= 0; i-- {
			lgk := Lagerklassen[i]
			if strings.Contains(match, lgk) {
				e.result.Lagerklasse = lgk
				return nil
			}
		}
	}
	return errors.New("Could not extract lagerklasse.")
}

func (e *DefaultExtractor) ExtractHPSaetze() error {
	matches := hpSatzRegex.FindAllString(e.content, -1)

	if len(matches) == 0 {
		return errors.New("Could not extract HP-Sätze.")
	}

	for _, hpsatz := range matches {
		hpsatz = strings.TrimSpace(hpsatz)
		if strings.Contains(hpsatz, "H") {
			e.result.HSaetze = append(e.result.HSaetze, hpsatz)
		} else {
			e.result.PSaetze = append(e.result.PSaetze, hpsatz)
		}
	}
	e.result.HSaetze = helpers.RemoveDuplicates(e.result.HSaetze)
	e.result.PSaetze = helpers.RemoveDuplicates(e.result.PSaetze)
	sort.Slice(e.result.HSaetze, func(i, j int) bool { return e.result.HSaetze[i] < e.result.HSaetze[j] })
	sort.Slice(e.result.PSaetze, func(i, j int) bool { return e.result.PSaetze[i] < e.result.PSaetze[j] })
	return nil
}

func (e *DefaultExtractor) Extract() *SicherheitsdatenblattData {
	// TODO: Execute in Parallel

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		if err := e.ExtractBezeichnung(); err != nil {
			e.result.Bezeichnung = err.Error()
			hclog.Default().Error("error", "error", err)
		}
		wg.Done()
	}()

	go func() {
		if err := e.ExtractSignalwort(); err != nil {
			e.result.Signalwort = err.Error()
			hclog.Default().Error("error", "error", err)
		}
		wg.Done()
	}()

	go func() {
		if err := e.ExtractLagerklasse(); err != nil {
			e.result.Lagerklasse = err.Error()
			hclog.Default().Error("error", "error", err)
		}
		wg.Done()
	}()

	go func() {
		if err := e.ExtractHPSaetze(); err != nil {
			hclog.Default().Error("error", "error", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return e.result
}
