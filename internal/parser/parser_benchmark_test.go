package parser_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/devbackend/crutch-db/internal/parser"
)

func BenchmarkParser_ParseSimpleValues(b *testing.B) {
	benchmarks := []struct {
		name    string
		command string
	}{
		{
			name:    "simple GET",
			command: "GET key-name",
		},
		{
			name:    "SET int",
			command: "SET key-name 123",
		},
		{
			name:    "SET string",
			command: "SET key-name 'hello'",
		},
		{
			name:    "SET array of int numbers",
			command: "SET key-name [1, 2, 3]",
		},
		{
			name:    "SET array of strings",
			command: "SET key-name ['hello', 'world']",
		},
		{
			name:    "simple DELETE",
			command: "DELETE key-name",
		},
		{
			name:    "simple KEYS",
			command: "KEYS",
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			p := parser.New()

			for i := 0; i < b.N; i++ {
				_, _ = p.Parse(bm.command)
			}
		})
	}
}

func BenchmarkParser_ParseSetBigValues(b *testing.B) {
	benchmarks := []struct {
		name  string
		value func() string
	}{
		{
			name: "string 1kb",
			value: func() string {
				chars := make([]byte, 1024)
				for k := range chars {
					chars[k] = 'a'
				}

				return fmt.Sprintf("'%s'", string(chars))
			},
		},
		{
			name: "string 1mb",
			value: func() string {
				chars := make([]byte, 1024*1024)
				for k := range chars {
					chars[k] = 'a'
				}

				return fmt.Sprintf("'%s'", string(chars))
			},
		},
		{
			name: "array of 1K int numbers",
			value: func() string {
				numbers := make([]string, 1000)
				for k := range numbers {
					numbers[k] = strconv.Itoa(k)
				}

				return fmt.Sprintf("[%s]", strings.Join(numbers, ","))
			},
		},
		{
			name: "array of 1M int numbers",
			value: func() string {
				numbers := make([]string, 1_000_000)
				for k := range numbers {
					numbers[k] = strconv.Itoa(k)
				}

				return fmt.Sprintf("[%s]", strings.Join(numbers, ","))
			},
		},
		{
			name: "array of 1K strings",
			value: func() string {
				words := make([]string, 1000)
				for k := range words {
					words[k] = fmt.Sprintf("'word-%s'", strconv.Itoa(k))
				}

				return fmt.Sprintf("[%s]", strings.Join(words, ","))
			},
		},
		{
			name: "array of 1M strings",
			value: func() string {
				words := make([]string, 1_000_000)
				for k := range words {
					words[k] = fmt.Sprintf("'word-%s'", strconv.Itoa(k))
				}

				return fmt.Sprintf("[%s]", strings.Join(words, ","))
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			p := parser.New()

			cmd := "SET key-name " + bm.value()

			for i := 0; i < b.N; i++ {
				_, _ = p.Parse(cmd)
			}
		})
	}
}
