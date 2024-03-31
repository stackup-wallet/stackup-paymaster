install-dev:
	go install github.com/cosmtrek/air@latest
	go mod tidy
	cp .env.template .env

dev-run:
	air -c .air.server.toml

generate-contract-pkg:
	abigen --abi=./abi/verifyingPaymaster.json --pkg=contract --out=./pkg/contract/bindings.go