all: bench test
bench: 
	@go test -bench=. -benchtime=3s -benchmem -benchmem github.com/senayuki/mosaic/processor
	@go test -bench=. -benchtime=3s -benchmem -benchmem github.com/senayuki/mosaic/mask
test: 
	@go test -covermode=count -coverprofile=processor.cov -timeout 30s ./...