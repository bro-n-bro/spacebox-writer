package clickhouse

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func BenchmarkMarshalToString(b *testing.B) {
	data := []byte(`{"a": "1", "b": 2, "c": 3}`)

	b.Run("byte to string", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bytes, err := jsoniter.Marshal(data)
			if err != nil {
				b.Fatal(err)
			}
			_ = string(bytes)
		}
	})

	b.Run("marshal string", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsoniter.MarshalToString(data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
