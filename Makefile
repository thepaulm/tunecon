builddir := build

$(info builddir ${builddir})

all: bins

${builddir}:
	mkdir -p $(builddir)

.PHONY: bins
bins:
	go build -o ${builddir}/tuncon github.com/thepaulm/tunecon/cli

clean:
	rm -f ${builddir}/*
