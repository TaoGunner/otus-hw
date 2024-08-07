//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r, err := zip.OpenReader("testdata/users.dat.zip")
		if err != nil {
			panic(err)
		}

		data, err := r.File[0].Open()

		b.StartTimer()
		_, err = GetDomainStat(data, "com")
		b.StopTimer()
		if err != nil {
			panic(err)
		}
		r.Close()
	}
}
