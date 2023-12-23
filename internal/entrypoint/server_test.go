package entrypoint

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerateSeverEntryPoint(t *testing.T) {
	msg, _ := json.Marshal(map[string]string{"message": "message"})
	testCases := []struct {
		desc     string
		imports  []string
		filePath string
		props    json.RawMessage
		expect   string
	}{
		{
			desc:     "should add props in entry point",
			imports:  []string{},
			filePath: "test.jsx",
			props:    msg,
			expect: `const props = {"message":"message"}`,
		},
		{
			desc:     "should add imports in entry point",
			imports:  []string{`import { testApp } from "test";`},
			filePath: "test.jsx",
			props:    msg,
			expect: `import { testApp } from "test";`,
		},
		{
			desc:     "should add imports in entry point",
			imports:  []string{`import { testApp } from "test";`},
			filePath: "../test/test2/test.jsx",
			props:    msg,
			expect: `import App from "../test/test2/test.jsx";`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result, err := GenerateSeverEntryPoint(tC.imports, tC.filePath, tC.props)
			if err != nil {
				t.Fatalf("failed when expecting #%s \n err: %v", tC.expect, err)
			}
			if !strings.Contains(result, tC.expect) {
        t.Fatalf("\n\t expecting: %s \n\t to contains: \n %s \n", result, tC.expect)
			}
		})
	}
}
