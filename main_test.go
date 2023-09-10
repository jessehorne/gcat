package main

import (
	"os"
	"testing"
)

func Benchmark_getOptions(b *testing.B) {
	args := "-bvET"
	for i := 0; i < b.N; i++ {
		_, err := getOptions(args)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_parseArgs(b *testing.B) {
	args := []string{"-bvET", "-A", "-nAb", "test/1.txt"}
	for i := 0; i < b.N; i++ {
		_, _, err := parseArgs(args)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_gcat(b *testing.B) {
	files := []string{"test/1.txt"}

	var options Options
	options.OA = true
	options.Ob = true

	oldOut := os.Stdout
	os.Stdout = nil
	for i := 0; i < b.N; i++ {
		err := gcat(files, options)
		if err != nil {
			b.Error(err)
		}
	}
	os.Stdout = oldOut
}
