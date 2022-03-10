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
	Extract() *SicherheitsdatenblattData
}

type SicherheitsdatenblattData struct {
	Bezeichnung string   `json:"bezeichnung"`
	Lagerklasse string   `json:"lagerklasse"`
	Signalwort  string   `json:"signalwort"`
	HSaetze     []string `json:"hSaezte"`
	PSaetze     []string `json:"pSaezte"`
	GHS         []string `json:"ghs"`
	WGK         string   `json:"wgk"`
}

var (
	lagerklassen = [...]string{"1", "2A", "2B", "3", "4.1A", "4.1B", "4.2", "4.3", "5.1A", "5.1B", "5.1C", "5.2", "6.1A",
		"6.1B", "6.1C", "6.1D", "6.2", "7", "8A", "8B", "10", "11", "12", "13", "10-13"}
	hpSatzRegex       = regexp.MustCompile(`(?im)(\s?\+?\s?E?U?[HP][0-9]{3}[a-zA-Z]{0,2}){1,3}`)
	signalwortRegex   = regexp.MustCompile(`(?im)(signalwort|signalwörter)\r?\n?(.*)`)
	lagerklassenRegex = regexp.MustCompile(`(?im)lagerklasse(.*)`)
	bzRegex           = regexp.MustCompile(`(?im)((handels?)?name|produktidentifikator)(\s*)\n(.*)`)
	ghsRegex          = regexp.MustCompile(`(?im)ghs\s?-?[0-9]{2}`)
	wgkRegex          = regexp.MustCompile(`(?im)(Wassergefährdungsklasse|WGK)\s+?(\d)`)
)

type defaultExtractor struct {
	content string
	logger  hclog.Logger
}

func (e *defaultExtractor) ExtractBezeichnung(result *SicherheitsdatenblattData) error {
	matches := bzRegex.FindAllStringSubmatch(e.content, -1)

	if len(matches) == 0 {
		return errors.New("could not extract bezeichnung")
	}

	for i := range matches {
		for j := 1; j < len(matches[i]); j++ {
			candidate := matches[i][j]

			candidate = strings.ToLower(candidate)

			// Strip Keyword
			candidate = strings.Replace(candidate, "produktidentifikator", "", -1)
			candidate = strings.Replace(candidate, "produktname", "", -1)
			candidate = strings.Replace(candidate, "handelsname", "", -1)
			candidate = strings.Replace(candidate, ":", "", -1)

			// Trim
			candidate = strings.TrimSpace(candidate)

			// Überschriftszeile
			if candidate == "" {
				continue
			}
			result.Bezeichnung = candidate
			return nil
		}
	}

	return errors.New("could not extract bezeichnung")
}

func (e *defaultExtractor) ExtractSignalwort(result *SicherheitsdatenblattData) error {
	matches := signalwortRegex.FindAllString(e.content, -1)

	if len(matches) == 0 {
		return errors.New("could not extract signalwort")
	}

	for _, match := range matches {
		match = strings.ToLower(match)
		if strings.Contains(match, "gefahr") {
			result.Signalwort = "Gefahr"
			return nil
		}
		if strings.Contains(match, "achtung") {
			result.Signalwort = "Achtung"
			return nil
		}
	}
	return errors.New("could not extract signalwort")
}

func (e *defaultExtractor) ExtractLagerklasse(result *SicherheitsdatenblattData) error {
	matches := lagerklassenRegex.FindAllString(e.content, -1)

	if len(matches) == 0 {
		return errors.New("could not extract lagerklasse")
	}

	for _, match := range matches {
		if strings.Contains(match, "510") {
			match = strings.ReplaceAll(match, "510", "")
		}

		match = strings.ReplaceAll(match, " ", "")

		// Um Highest-Value Matches abzufangen wird von hinten durchiteriert, damit 1 nicht alle 1x Werte abfängt
		for i := len(lagerklassen) - 1; i >= 0; i-- {
			lgk := lagerklassen[i]
			if strings.Contains(match, lgk) {
				result.Lagerklasse = lgk
				return nil
			}
		}
	}
	return errors.New("could not extract lagerklasse")
}

func (e *defaultExtractor) ExtractHPSaetze(result *SicherheitsdatenblattData) error {
	matches := hpSatzRegex.FindAllString(e.content, -1)

	if len(matches) == 0 {
		return errors.New("could not extract HP-Sätze")
	}

	for _, hpsatz := range matches {
		hpsatz = strings.TrimSpace(hpsatz)
		if strings.Contains(hpsatz, "H") {
			result.HSaetze = append(result.HSaetze, hpsatz)
		} else {
			result.PSaetze = append(result.PSaetze, hpsatz)
		}
	}
	result.HSaetze = helpers.RemoveDuplicates(result.HSaetze)
	result.PSaetze = helpers.RemoveDuplicates(result.PSaetze)

	sort.Strings(result.HSaetze)
	sort.Strings(result.PSaetze)

	return nil
}

func (e *defaultExtractor) ExtractGHS(result *SicherheitsdatenblattData) error {
	matches := ghsRegex.FindAllString(e.content, -1)

	if len(matches) == 0 {
		return errors.New("could not extract GHS")
	}

	for _, match := range matches {
		match = strings.TrimSpace(match)
		result.GHS = append(result.GHS, match)
	}

	return nil
}

func (e *defaultExtractor) ExtractWGK(result *SicherheitsdatenblattData) error {
	matches := wgkRegex.FindAllStringSubmatch(e.content, -1)

	if len(matches) == 0 {
		return errors.New("could not extract WGK")
	}

	for _, match := range matches {
		wgk := match[2] // second capturing groups contains the number
		result.WGK = wgk
	}

	return nil
}

func (e *defaultExtractor) Extract() *SicherheitsdatenblattData {
	var wg sync.WaitGroup
	wg.Add(6)

	result := &SicherheitsdatenblattData{}

	go func() {
		if err := e.ExtractBezeichnung(result); err != nil {
			result.Bezeichnung = err.Error()
			e.logger.Error("error", "error", err)
		}
		wg.Done()
	}()

	go func() {
		if err := e.ExtractSignalwort(result); err != nil {
			result.Signalwort = err.Error()
			e.logger.Error("error", "error", err)
		}
		wg.Done()
	}()

	go func() {
		if err := e.ExtractLagerklasse(result); err != nil {
			result.Lagerklasse = err.Error()
			e.logger.Error("error", "error", err)
		}
		wg.Done()
	}()

	go func() {
		if err := e.ExtractHPSaetze(result); err != nil {
			e.logger.Error("error", "error", err)
		}
		wg.Done()
	}()

	go func() {
		if err := e.ExtractGHS(result); err != nil {
			e.logger.Error("error", "error", err)
		}
		wg.Done()
	}()

	go func() {
		if err := e.ExtractWGK(result); err != nil {
			e.logger.Error("error", "error", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return result
}

func Extract(content string, logger hclog.Logger) *SicherheitsdatenblattData {
	e := &defaultExtractor{content: content, logger: logger}
	return e.Extract()
}
