-include .env.develop

START_LOG = @echo "================================================= START OF LOG ==================================================="
END_LOG = @echo "================================================== END OF LOG ===================================================="

.PHONY: env
env: ./.env.develop
	$(START_LOG)
	@cp ./.env.develop.tmpl ./.env.develop
	@touch .cartesi.env
	@echo "Environment file created at ./.env.develop"
	$(END_LOG)

.PHONY: build
build:
	$(START_LOG)
	@docker build \
		-t machine:latest \
		-f ./build/Dockerfile.machine .
	@cartesi build --from-image machine:latest
	$(END_LOG)

.PHONY: dev
dev:
	$(START_LOG)
	@cd ./tools/tlsnotary/verifier && cargo build --release
	@cp ./tools/tlsnotary/verifier/target/release/libverifier.a ./internal/usecase/crowdfunding_usecase/
	@nonodo -- air
	
.PHONY: generate
generate:
	$(START_LOG)
	@go run ./pkg/rollups-contracts/generate
	$(END_LOG)

.PHONY: test
test:
	@cd ./tools/tlsnotary/verifier && cargo build --release
	@cp ./tools/tlsnotary/verifier/target/release/libverifier.a ./internal/usecase/crowdfunding_usecase/
	@go test -p=1 ./... -coverprofile=./coverage.md -v

.PHONY: coverage
coverage: test
	@go tool cover -html=./coverage.md

.PHONY: state
state:
	@chmod +x ./tools/state.sh
	@./tools/state.sh