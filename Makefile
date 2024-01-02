pprof:
	go build -o toolbox -tags "pprof" main.go

all:
	go build -o toolbox -tags "" main.go`