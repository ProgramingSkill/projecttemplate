SRC = $(wildcard *.go)
GITVISION = `git describe`"-"`git rev-parse --short HEAD`

all:projecttemplate
projecttemplate:$(SRC)
	go build -x -ldflags "-w -s -X main.BuildVersion=$(GITVISION)" -o projecttemplate $^
clean:
	rm -rvf projecttemplate
