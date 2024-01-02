pprof:
	go build -o toolbox -tags="pprof" main.go

jsoncheck:
	go build -o toolbox -tags="jsoncheck" main.go

all:
	go build -o toolbox -tags="pprof,jsoncheck" main.go