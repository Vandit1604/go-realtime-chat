.PHONY: postgresinit
postgresinit:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -e POSTGRES_DB=go-chat -d postgres:15-alpine

.PHONY: postgres
postgres:
	docker exec -it postgres psql go-chat

.PHONY: createdb
createdb:
	docker exec -it postgres createdb --username=root --owner=root go-chat

.PHONY: dropdb
dropdb:
	docker exec -it postgres dropdb go-chat
