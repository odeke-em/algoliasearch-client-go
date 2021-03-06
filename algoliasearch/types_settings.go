package algoliasearch

import (
	"fmt"
	"os"
	"reflect"
)

// Settings is the structure returned by `GetSettigs` to ease the use of the
// index settings.
type Settings struct {
	// Indexing parameters
	AllowCompressionOfIntegerArray bool     `json:"allowCompressionOfIntegerArray"`
	AttributeForDistinct           string   `json:"attributeForDistinct"`
	AttributesForFaceting          []string `json:"attributesForFaceting"`
	AttributesToIndex              []string `json:"attributesToIndex"`
	CustomRanking                  []string `json:"customRanking"`
	NumericAttributesToIndex       []string `json:"numericAttributesToIndex"`
	NumericAttributesForFiltering  []string `json:"numericAttributesForFiltering"`
	Ranking                        []string `json:"ranking"`
	Replicas                       []string `json:"replicas"`
	SearchableAttributes           []string `json:"searchableAttributes"`
	SeparatorsToIndex              string   `json:"separatorsToIndex"`
	Slaves                         []string `json:"slaves"`
	UnretrievableAttributes        []string `json:"unretrievableAttributes"`

	// Query expansion
	DisableTypoToleranceOnAttributes []string `json:"disableTypoToleranceOnAttributes"`
	DisableTypoToleranceOnWords      []string `json:"disableTypoToleranceOnWords"`

	// Default query parameters (can be overridden at query-time)
	AdvancedSyntax             bool        `json:"advancedSyntax"`
	AllowTyposOnNumericTokens  bool        `json:"allowTyposOnNumericTokens"`
	AttributesToHighlight      []string    `json:"attributesToHighlight"`
	AttributesToRetrieve       []string    `json:"attributesToRetrieve"`
	AttributesToSnippet        []string    `json:"attributesToSnippet"`
	Distinct                   interface{} `json:"distinct"` // float64 (actually an int) or bool
	HighlightPostTag           string      `json:"highlightPostTag"`
	HighlightPreTag            string      `json:"highlightPreTag"`
	HitsPerPage                int         `json:"hitsPerPage"`
	IgnorePlurals              bool        `json:"ignorePlurals"`
	MaxValuesPerFacet          int         `json:"maxValuesPerFacet"`
	MinProximity               int         `json:"minProximity"`
	MinWordSizefor1Typo        int         `json:"minWordSizefor1Typo"`
	MinWordSizefor2Typos       int         `json:"minWordSizefor2Typos"`
	OptionalWords              []string    `json:"optionalWords"`
	QueryType                  string      `json:"queryType"`
	RemoveStopWords            interface{} `json:"removeStopWords"` // []interface{} (actually a []string) or bool
	ReplaceSynonymsInHighlight bool        `json:"replaceSynonymsInHighlight"`
	SnippetEllipsisText        string      `json:"snippetEllipsisText"`
	TypoTolerance              string      `json:"typoTolerance"`
}

// clean sets the nil `interface{}` fields of any `Settings struct` generated
// by `GetSettings`.
func (s *Settings) clean() {
	if s.Distinct == nil {
		s.Distinct = false
	}

	if s.RemoveStopWords == nil {
		s.RemoveStopWords = false
	}
}

// ToMap produces a `Map` corresponding to the `Settings struct`. It should
// only be used when it's needed to pass a `Settings struct` to `SetSettings`,
// typically when one needs to copy settings between two indices.
func (s *Settings) ToMap() Map {
	// Add all fields except:
	//  - RemoveStopWords []interface{} or bool
	//  - Distinct float64 or bool
	//  - TypoTolerance string or bool
	m := Map{
		// Indexing parameters
		"allowCompressionOfIntegerArray": s.AllowCompressionOfIntegerArray,
		"attributeForDistinct":           s.AttributeForDistinct,
		"attributesForFaceting":          s.AttributesForFaceting,
		"attributesToIndex":              s.AttributesToIndex,
		"customRanking":                  s.CustomRanking,
		"numericAttributesToIndex":       s.NumericAttributesToIndex,
		"numericAttributesForFiltering":  s.NumericAttributesForFiltering,
		"ranking":                        s.Ranking,
		"replicas":                       s.Replicas,
		"searchableAttributes":           s.SearchableAttributes,
		"separatorsToIndex":              s.SeparatorsToIndex,
		"slaves":                         s.Slaves,
		"unretrievableAttributes":        s.UnretrievableAttributes,

		// Query expansion
		"disableTypoToleranceOnAttributes": s.DisableTypoToleranceOnAttributes,
		"disableTypoToleranceOnWords":      s.DisableTypoToleranceOnWords,

		// Default query parameters (can be overridden at query-time)
		"advancedSyntax":             s.AdvancedSyntax,
		"allowTyposOnNumericTokens":  s.AllowTyposOnNumericTokens,
		"attributesToHighlight":      s.AttributesToHighlight,
		"attributesToRetrieve":       s.AttributesToRetrieve,
		"attributesToSnippet":        s.AttributesToSnippet,
		"highlightPostTag":           s.HighlightPostTag,
		"highlightPreTag":            s.HighlightPreTag,
		"hitsPerPage":                s.HitsPerPage,
		"ignorePlurals":              s.IgnorePlurals,
		"maxValuesPerFacet":          s.MaxValuesPerFacet,
		"minProximity":               s.MinProximity,
		"minWordSizefor1Typo":        s.MinWordSizefor1Typo,
		"minWordSizefor2Typos":       s.MinWordSizefor2Typos,
		"optionalWords":              s.OptionalWords,
		"queryType":                  s.QueryType,
		"replaceSynonymsInHighlight": s.ReplaceSynonymsInHighlight,
		"snippetEllipsisText":        s.SnippetEllipsisText,
	}

	// Remove empty string slices to avoid creating null-valued fields in the
	// JSON settings sent to the API
	var sliceAttributesToRemove []string

	for attr, value := range m {
		switch v := value.(type) {
		case []string:
			if len(v) == 0 {
				sliceAttributesToRemove = append(sliceAttributesToRemove, attr)
			}
		}
	}

	for _, attr := range sliceAttributesToRemove {
		delete(m, attr)
	}

	// Handle `RemoveStopWords` separately as it may be either a `bool` or a
	// `[]interface{}` which is in fact a `[]string`.
	switch v := s.RemoveStopWords.(type) {

	case bool:
		m["removeStopWords"] = v

	case []interface{}:
		var languages []string
		for _, itf := range v {
			lang, ok := itf.(string)
			if ok {
				languages = append(languages, lang)
			} else {
				fmt.Fprintln(os.Stderr, "Settings.ToMap(): `removeStopWords` slice doesn't only contain strings")
			}
		}
		if len(languages) > 0 {
			m["removeStopWords"] = languages
		}

	default:
		fmt.Println(reflect.TypeOf(s.RemoveStopWords))
		fmt.Fprintln(os.Stderr, "Settings.ToMap(): Wrong type for `removeStopWords`")

	}

	// Handle `Distinct` separately as it may be either a `bool` or a `float64`
	// which is in fact a `int`.
	switch v := s.Distinct.(type) {
	case bool:
		m["distinct"] = v
	case float64:
		m["distinct"] = int(v)
	}

	return m
}
