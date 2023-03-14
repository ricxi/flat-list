include .envrc
export

.PHONY: run/dev/token
run/dev/token:
	go run .

.PHONY: run/dev/psql/token
run/dev/psql/token:
	@echo "starting postgres for token service"
	@./postgres.sh ${PSQL_DSN}