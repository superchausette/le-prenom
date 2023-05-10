package leprenom

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestImport(t *testing.T) {
	tests := []struct {
		Name   string
		Input  string
		Output []CsvContent
		Error  error
	}{
		{
			Name:   "Empty",
			Input:  "",
			Output: []CsvContent{},
		},
		{
			Name:   "OneLineBoy",
			Input:  "1;AADAM;2009;4",
			Output: []CsvContent{{1, "AADAM", 2009, 4}},
		},
		{
			Name:   "OneLineGirl",
			Input:  "2;CÉLIA;2017;707",
			Output: []CsvContent{{2, "CÉLIA", 2017, 707}},
		},
		{
			Name:  "MultipleLineBoysAndGirls",
			Input: "1;AADAM;2009;4\n2;CÉLIA;2017;707\n2;CELIANE;1930;6\n2;CELIANE;1934;7",
			Output: []CsvContent{
				{1, "AADAM", 2009, 4},
				{2, "CÉLIA", 2017, 707},
				{2, "CELIANE", 1930, 6},
				{2, "CELIANE", 1934, 7},
			},
		},
		{
			Name:   "OneLineRareFilter",
			Input:  "2;_PRENOMS_RARES;2017;707",
			Output: []CsvContent{},
		},
		{
			Name:   "OneLineGirlMissingColumn",
			Input:  "2;CÉLIA;2017",
			Output: []CsvContent{},
			Error:  errors.New("Invalid size of record [2 CÉLIA 2017] 3, expected 4"),
		},
		{
			Name:   "OneLineGirlExtraColumn",
			Input:  "2;CÉLIA;2017;797;Hey",
			Output: []CsvContent{},
			Error:  errors.New("Invalid size of record [2 CÉLIA 2017 797 Hey] 5, expected 4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			output, err := Import(strings.NewReader(tt.Input))
			if !reflect.DeepEqual(err, tt.Error) {
				t.Fatalf("Import(%s) unexpected error got: '%s', want: '%s'", tt.Input, err, tt.Error)
			}
			if !reflect.DeepEqual(output, tt.Output) {
				t.Fatalf("Import(%s) = %v, want %v", tt.Input, output, tt.Output)
			}
		})
	}
}
