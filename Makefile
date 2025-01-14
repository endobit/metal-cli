BUILDER=./.builder
RULES=go
include $(BUILDER)/rules.mk
$(BUILDER)/rules.mk:
	-go run github.com/endobit/builder@latest init

GO = gotip
BIN = bin

format::
	buf format -w

lint::
#	cd proto && buf lint
	gotip tool github.com/sqlc-dev/sqlc/cmd/sqlc compile

generate::
	buf generate
	gotip tool github.com/sqlc-dev/sqlc/cmd/sqlc generate

./$(BIN):
	mkdir $@

build:: ./$(BIN)
	$(GO_BUILD) -o $(BIN) ./cmd/stack
	$(GO_BUILD) -o $(BIN) ./cmd/stackd

clean::
	rm -rf $(BIN)

nuke::
	rm -rf internal/generated



