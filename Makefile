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

.PHONY: machine
machine:
	$(START_LOG)
	@docker build \
		-t machine:latest \
		-f ./build/Dockerfile.machine .
	@cartesi build --from-image machine:latest
	$(END_LOG)

.PHONY: local
local:
	$(START_LOG)
	@nonodo -- cartesi-machine --network \
		--flash-drive=label:root,filename:.cartesi/image.ext2 \
		--env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004 \
		-- /var/opt/cartesi-app/app
	
.PHONY: bidings
bidings:
	$(START_LOG)
	@go run ./pkg/rollups_contracts/generate
	$(END_LOG)

.PHONY: test
test:
	@go test -p 1 ./... -coverprofile=./coverage.md -v

.PHONY: deploy_token
deploy_token:
	$(START_LOG)
	@cd contracts && forge script script/DeployVoltToken.s.sol --rpc-url $(RPC_URL) --broadcast --verify --etherscan-api-key $(API_KEY) -vvv
	$(END_LOG)

.PHONY: coverage
coverage: test
	@go tool cover -html=./coverage.md

.PHONY: docs
docs:
	@cd docs && npm run dev