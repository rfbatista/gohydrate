package gohydrate

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	basepath, _ := os.Getwd()
	testCases := []struct {
		desc       string
		pageConfig PageConfig
		expectHTML string
	}{
		{
			desc:       "if page is rendered",
			pageConfig: PageConfig{Filename: "app.jsx", Props: map[string]string{"mensage": "test mensagem"}},
			expectHTML: "<h1>Hello from React!</h1>",
		},
		{
			desc:       "teste 2",
			pageConfig: PageConfig{Filename: "app.jsx", Props: map[string]string{"mensage": "test mensagem"}},
			expectHTML: "<h1>Hello from React!</h1>",
		},
	}
	for _, tC := range testCases {
		e, _ := New(EngineConfig{
			BasePath:  basepath,
			PagesPath: "/examples/ssr",
		})
		t.Run(tC.desc, func(t *testing.T) {
			r, err := e.RenderPage(tC.pageConfig)
			if err != nil {
				t.Fatalf("failed when expecting #%s \n err: %v", tC.expectHTML, err)
			}
			if r.HTML != tC.expectHTML {
				t.Fatalf("\n expected: %s \n receive: %s", tC.expectHTML, r.HTML)
			}
		})
	}
}

func BenchmarkEngine(b *testing.B) {
	basepath, _ := os.Getwd()
	var e, _ = New(EngineConfig{
		BasePath:  basepath,
		PagesPath: "/examples/ssr",
	})

	testCases := []struct {
		desc       string
		pageConfig PageConfig
		expectHTML string
	}{
		{
			desc:       "if page is rendered",
			pageConfig: PageConfig{Filename: "app.jsx", Props: map[string]string{"mensage": "test mensagem"}},
			expectHTML: "<h1>Hello from React!</h1>",
		},
	}
  b.ResetTimer()
  b.StartTimer()
	for i := 0; i < b.N; i++ {
		tC := testCases[0]
		r, err := e.RenderPage(tC.pageConfig)
		if err != nil {
			b.Fatalf("failed when expecting #%s \n err: %v", tC.expectHTML, err)
		}
		if r.HTML != tC.expectHTML {
			b.Fatalf("\n expected: %s \n receive: %s", tC.expectHTML, r.HTML)
		}
	}
  b.StopTimer()
	fmt.Println("\nNumber of iterations: ", b.N)
	fmt.Println("Elapsed:", b.Elapsed()/time.Duration(b.N))
}
