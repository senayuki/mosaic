all: bench test
bench: 
	@go test -bench=. -benchtime=3s -benchmem github.com/senayuki/mosaic/processor
	@go test -bench=. -benchtime=3s -benchmem github.com/senayuki/mosaic/mask
test: 
	@go test -covermode=count -coverprofile=processor.cov -timeout 30s github.com/senayuki/mosaic/processor
	@go test -covermode=count -coverprofile=mask.cov -timeout 30s github.com/senayuki/mosaic/mask