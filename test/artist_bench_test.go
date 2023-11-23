// test/artist_bench_test.go

package test

import (
	"github.com/blnto/blnto_service/internal/domain/artist"
	"testing"
)

func BenchmarkDirectCopy(b *testing.B) {
	for n := 0; n < b.N; n++ {
		original := CreateMockArtist()
		copy := *original
		_ = &copy
	}
}

func BenchmarkCopyWithModification(b *testing.B) {
	for n := 0; n < b.N; n++ {
		original := CreateMockArtist()
		copy := *original
		modifiedCopy := &copy
		modifiedCopy.Name = "New Name"
		_ = modifiedCopy
	}
}

func BenchmarkCreateMockArtist(b *testing.B) {
	original := CreateMockArtist()
	NewArtistCopy(original, func(artist *artist.Artist) {
		original.Name = "New Name"
		original.Location = "New Location"
	})
}

// Additional benchmarks can be added here...
