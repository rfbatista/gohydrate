package gohydrate

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	basepath, _ := os.Getwd()
	e, _ := New(EngineConfig{
		BasePath:  basepath,
		PagesPath: "/examples/ssr",
	})
	t.Run("should render page with props", func(t *testing.T) {
		expect := "<h1>Hello from React!<div>teste</div></h1>"
		c := PageConfig{Filename: "app.jsx", Props: map[string]string{"mensage": "test mensagem", "title": "teste"}}
		r, err := e.RenderPage(c)
		if err != nil {
			t.Fatalf("failed")
		}
		if r.HTML != expect {
			t.Fatalf("\n expected: %s \n receive: %s", expect, r.HTML)
		}
	})
	t.Run("should render page without props", func(t *testing.T) {
		expect := "<h1>Hello from React!<div></div></h1>"
		c := PageConfig{Filename: "app.jsx", Props: nil}
		r, err := e.RenderPage(c)
		if err != nil {
			t.Fatalf("failed")
		}
		if r.HTML != expect {
			t.Fatalf("\n expected: %s \n receive: %s", expect, r.HTML)
		}
	})
	t.Run("should fail if file is not found", func(t *testing.T) {
		c := PageConfig{Filename: "test.jsx", Props: map[string]string{"mensage": "test mensagem", "title": "teste 2"}}
		_, err := e.RenderPage(c)
		if err == nil {
			t.Fatalf("should have failed if file is not found")
		}
	})
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
			pageConfig: PageConfig{Filename: "app.jsx", Props: map[string]string{"title": "test"}},
			expectHTML: "<h1>Hello from React!<div>test</div></h1>",
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
		_ = r
	}
	b.StopTimer()
	fmt.Println("\nNumber of iterations: ", b.N)
	fmt.Println("Elapsed:", b.Elapsed()/time.Duration(b.N))
}
