guard-%:
	@#$(or ${$*}, $(error Environment variable $* is not set))

test: db.test.prepare
	go test -v ./...

db.test.prepare: guard-TEST_DATABASE guard-TEST_DATABASE_URL
	createdb ${TEST_DATABASE} || true
	env GOOSE_DBSTRING=${TEST_DATABASE_URL} goose up
	psql ${TEST_DATABASE} -f db/test_seeds.sql
