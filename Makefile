all: bench

test: *.go
	go test -race .

bench: *.go
	go test -coverprofile=coverage.out .
	go tool cover -func=coverage.out
	go tool cover -html=coverage.html

	go test -covermode=count -coverprofile=count.out .
	go tool cover -func=count.out
	go tool cover -html=count.html

	# no bugs
	go test -race -run=none -bench=BenchmarkGet .

	# cpu profile
	go test -run=none -bench=BenchmarkGet -cpuprofile=cprof .
	go tool pprof --text lru.test cprof

	# memory allocations
	go test -run=none -bench=BenchmarkGet -memprofile=mprof -memprofilerate=1 .
	go tool pprof --alloc_space --text lru.test mprof
	go tool pprof --alloc_objects --text lru.test mprof
	go tool pprof --inuse_space --text lru.test mprof
	go tool pprof --inuse_objects --text lru.test mprof

	# blocking profile
	go test -run=none -bench=BenchmarkGet -blockprofile=blockprof -blockprofilerate=1 .
	go tool pprof --text --lines lru.test blockprof

