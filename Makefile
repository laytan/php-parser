all: compile fmt build

fmt:
	find . -type f -iregex '.*\.go' -exec gofmt -l -s -w '{}' +

build:
	go generate ./...
	go build ./cmd/...

test:
	go test -v ./...

cover:
	go test ./... --cover

bench:
	go test -benchmem -bench=. ./internal/php8
	go test -benchmem -bench=. ./internal/php7

compile: ./internal/php8/php8.go ./internal/php7/php7.go ./internal/php8/scanner.go ./internal/php7/scanner.go
	sed -i -e 's/yyErrorVerbose = false/yyErrorVerbose = true/g' ./internal/php8/php8.go
	sed -i -e 's/yyErrorVerbose = false/yyErrorVerbose = true/g' ./internal/php7/php7.go
	sed -i -e 's/\/\/line/\/\/ line/g' ./internal/php8/php8.go
	sed -i -e 's/\/\/line/\/\/ line/g' ./internal/php7/php7.go
	sed -i -e 's/\/\/line/\/\/ line/g' ./internal/php8/scanner.go
	sed -i -e 's/\/\/line/\/\/ line/g' ./internal/php7/scanner.go
	rm -f y.output

./internal/php8/php8.go: ./internal/php8/php8.y
	goyacc -o $@ $<

./internal/php8/scanner.go: ./internal/php8/scanner.rl
	ragel -Z -G2 -o $@ $<

./internal/php7/scanner.go: ./internal/php7/scanner.rl
	ragel -Z -G2 -o $@ $<

./internal/php7/php7.go: ./internal/php7/php7.y
	goyacc -o $@ $<

cpu_pprof:
	go test -cpuprofile cpu.pprof -bench=. -benchtime=20s ./internal/php8
	go tool pprof ./php8.test cpu.pprof

mem_pprof:
	go test -memprofile mem.pprof -bench=. -benchtime=20s -benchmem ./internal/php8
	go tool pprof -alloc_objects ./php8.test mem.pprof

cpu_pprof_php7:
	go test -cpuprofile cpu.pprof -bench=. -benchtime=20s ./internal/php7
	go tool pprof ./php7.test cpu.pprof

mem_pprof_php5:
	go test -memprofile mem.pprof -bench=. -benchtime=20s -benchmem ./internal/php7
	go tool pprof -alloc_objects ./php7.test mem.pprof
