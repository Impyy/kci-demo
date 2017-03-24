all: prep
	go build -o build/bin/kci-demo github.com/Impyy/kci-demo/cmd/kci-demo

prep:
	mkdir -p build/bin

clean:
	rm -rf build
