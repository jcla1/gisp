package main

import "testing"

func BenchmarkFuncCallSpeedWithInner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		withInner()
	}
}

func BenchmarkFuncCallSpeedWithoutInner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		withInner()
	}
}

func withInner() int {
	return func() int {
		return 10
	}()
}

func withoutInner() int {
	return 10
}
