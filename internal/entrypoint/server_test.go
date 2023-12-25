package entrypoint

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerateSeverEntryPoint(t *testing.T) {
	testCases := []struct {
		desc     string
		imports  []string
		filePath string
		props    json.RawMessage
		expect   string
	}{
		{
			desc:     "should add imports in entry point",
			imports:  []string{`import { testApp } from "test";`},
			filePath: "test.jsx",
			expect: `import { testApp } from "test";`,
		},
		{
			desc:     "should add imports in entry point",
			imports:  []string{`import { testApp } from "test";`},
			filePath: "../test/test2/test.jsx",
			expect: `import App from "../test/test2/test.jsx";`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result, err := GenerateSeverEntryPoint(tC.imports, tC.filePath)
			if err != nil {
				t.Fatalf("failed when expecting #%s \n err: %v", tC.expect, err)
			}
			if !strings.Contains(result, tC.expect) {
        t.Fatalf("\n\t expecting: %s \n\t to contains: \n %s \n", result, tC.expect)
			}
		})
	}
}
