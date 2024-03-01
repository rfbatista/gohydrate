package gohydrate

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

func NewTailwind(log *zap.Logger) *Tailwind {
	return &Tailwind{Logger: log}
}

type Tailwind struct {
	Logger     *zap.Logger
	addconfig  bool
	config     string
	addpostcss bool
	postcss    string
}

func (t *Tailwind) AddConfig(path string) {
	t.addconfig = true
	t.config = path
}

func (t *Tailwind) AddPostCSSConfig(path string) {
	t.addpostcss = true
	t.postcss = path
}

func (t *Tailwind) Build(buildIn string) []byte {
	var params []string
	params = append(params, "tailwind")
	if t.addconfig {
		params = append(params, "-c")
		params = append(params, t.config)
	}
	if t.addpostcss {
		params = append(params, "--postcss")
		params = append(params, t.postcss)
	}
	cmd := exec.Command(
		"npx",
		params...,
	)
	t.Logger.Debug(cmd.String())
	cmd.Dir = buildIn
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		t.Logger.Error(
			"failed to create tailwind css",
			zap.Error(err),
			zap.String("stdout", errb.String()),
		)
		log.Fatal(err)
	}
	return outb.Bytes()
}

func (t *Tailwind) WriteTo(path string, from string, filename string) error {
	out := t.Build(from)
	err := os.WriteFile(path+"/"+filename, out, 0o644)
	if err != nil {
		return err
	}
	return nil
}
