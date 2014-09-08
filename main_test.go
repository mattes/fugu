package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var parseFlagsTests = []struct {
	in  []string
	out []string
}{
	{
		[]string{`-flag2`, `dieter`, `kranker`, `-b2`, `heinz`},
		[]string{"-flag", "-flag=x", `-flag="x"`, `-flag='x'`, "-flag x"},
	},
	// {
	// 	`-flag -flag=x -flag="x" -flag='x' -flag x a`,
	// 	[]string{"-flag", "-flag=x", `-flag="x"`, `-flag='x'`, "-flag x"},
	// },
	// {
	// 	`a -flag -flag=x -flag="x" -flag='x' -flag x`,
	// 	[]string{"-flag", "-flag=x", `-flag="x"`, `-flag='x'`, "-flag x"},
	// },
	// {
	// 	`-flag -flag=x a", -flag="x" -flag='x' -flag x`,
	// 	[]string{"-flag", "-flag=x", `-flag="x"`, `-flag='x'`, "-flag x"},
	// },
	// {
	// 	`-flag -flag=x -flag="x" -flag='x' -rm a`,
	// 	[]string{"-flag", "-flag=x", `-flag="x"`, `-flag='x'`, "-rm"},
	// },
	// {
	// 	`-flag -flag=x -flag="x" -flag='x' -rm=true a`,
	// 	[]string{"-flag", "-flag=x", `-flag="x"`, `-flag='x'`, "-rm=true"},
	// },
}

func TestParseFlags(t *testing.T) {
	for _, tt := range parseFlagsTests {
		out := parseFlags(tt.in)
		require.Equal(t, tt.out, out, fmt.Sprintf("%v", tt))
	}
}
