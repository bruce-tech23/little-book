.PHONY: docker
docker:
	@rm littlebook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -buildvcs=false -o littlebook .
	@docker rmi -f burcetech/littlebook:v0.0.1
	@docker build -t burcetech/littlebook:v0.0.1 .