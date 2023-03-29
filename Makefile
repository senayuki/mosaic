all: bench test
bench: 
	@go test -bench=. -benchtime=3s -benchmem github.com/senayuki/mosaic/processor
test: 
	@go test -timeout 30s github.com/senayuki/mosaic/processor