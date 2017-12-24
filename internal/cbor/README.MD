sprint @ cbor>pwd
/home/ravir/go/src/github.com/rs/zerolog/internal/cbor
sprint @ cbor>go test -v -benchmem -bench=.
=== RUN   TestAppendString
--- PASS: TestAppendString (0.00s)
=== RUN   TestAppendBytes
--- PASS: TestAppendBytes (0.00s)
=== RUN   TestAppendTimeNow
--- PASS: TestAppendTimeNow (0.00s)
=== RUN   TestAppendTimePastPresentInteger
--- PASS: TestAppendTimePastPresentInteger (0.00s)
=== RUN   TestAppendTimePastPresentFloat
--- PASS: TestAppendTimePastPresentFloat (0.00s)
=== RUN   TestAppendNull
--- PASS: TestAppendNull (0.00s)
=== RUN   TestAppendBool
--- PASS: TestAppendBool (0.00s)
=== RUN   TestAppendBoolArray
--- PASS: TestAppendBoolArray (0.00s)
=== RUN   TestAppendInt
--- PASS: TestAppendInt (0.00s)
=== RUN   TestAppendIntArray
--- PASS: TestAppendIntArray (0.00s)
=== RUN   TestAppendFloat32
--- PASS: TestAppendFloat32 (0.00s)
goos: linux
goarch: amd64
pkg: github.com/rs/zerolog/internal/cbor
BenchmarkAppendString/EncodingLast-4         	30000000	        56.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendString/MultiBytesFirst-4      	30000000	        47.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendString/MultiBytesMiddle-4     	30000000	        47.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendString/MultiBytesLast-4       	30000000	        47.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendString/NoEncoding-4           	30000000	        47.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendString/EncodingFirst-4        	30000000	        49.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendString/EncodingMiddle-4       	30000000	        52.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendTime/Integer-4                	30000000	        39.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendTime/Float-4                  	30000000	        56.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendInt/int-Positive-4            	100000000	        21.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendInt/int-Negative-4            	100000000	        20.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendInt/uint16-4                  	50000000	        30.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendInt/int8-4                    	100000000	        21.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendInt/int32-4                   	50000000	        26.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendInt/uint8-4                   	50000000	        29.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendInt/uint32-4                  	50000000	        37.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendInt/uint64-4                  	50000000	        36.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendInt/int16-4                   	50000000	        25.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendInt/int64-4                   	30000000	        39.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendFloat/Float32-4               	50000000	        23.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendFloat/Float64-4               	50000000	        33.1 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/rs/zerolog/internal/cbor	33.802s
