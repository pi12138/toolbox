pprof:
	go build -o toolbox -tags="pprof" main.go

jsoncheck:
	go build -o toolbox -tags="jsoncheck" main.go

webServer:
	go build -o toolbox -tags="webServer" main.go

all:
	go build -o toolbox -tags="pprof,jsoncheck,webServer" main.go