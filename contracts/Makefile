# LOADING ENV FILE
-include .env

# DEFAULT VARIABLES
START_LOG = @echo "==================== START OF LOG ===================="
END_LOG = @echo "==================== END OF LOG ======================"

define deploy
	$(START_LOG)
	@forge script script/Token.s.sol --private-key $(PRIVATE_KEY) --broadcast -vvv
	$(END_LOG)
endef

.PHONY: setup
setup: .env.tmpl
	cp .env.tmpl .env

.PHONY: token
token:
	@echo "Deploying Tokens..."
	@$(deploy)