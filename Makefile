test: | vendor
	go test github.com/factorysh/go-tawa/tawa

vendor:
	dep ensure
