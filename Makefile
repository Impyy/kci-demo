all: prep
	go build -o build/bin/kci-demo github.com/alexbakker/kci-demo/cmd/kci-demo

prep:
	mkdir -p build/bin

clean:
	rm -rf build
