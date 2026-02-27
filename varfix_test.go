package fixvars

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	var testdata = analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
}
