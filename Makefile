guard-%:
	@#$(or ${$*}, $(error Environment variable $* is not set))

test: db.test.prepare
	go test -count=1 -p=1 -parallel=1 -v ./...

db.test.prepare: guard-TEST_DATABASE guard-TEST_DATABASE_URL
	createdb ${TEST_DATABASE} || true
	env GOOSE_DBSTRING=${TEST_DATABASE_URL} goose up
	psql ${TEST_DATABASE} -f db/test_seeds.sql

db.create: guard-PGDATABASE
	createdb ${PGDATABASE} || true

db.seed:
	cd scripts && ./20250523_import_products.ps1

db.migrate: guard-GOOSE_DBSTRING
	goose up

db.setup: db.create db.migrate db.seed

db.drop: guard-PGDATABASE
	dropdb --if-exists ${PGDATABASE}

db.reset: db.drop db.setup
