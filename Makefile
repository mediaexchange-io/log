.PHONY: test bench

test:
	go test

bench:
	go test -cpu=1,2,4 -benchmem -benchtime=5s -bench "^Benchmark" -run=^$
