BUILDER=./.builder
RULES=go
include $(BUILDER)/rules.mk
$(BUILDER)/rules.mk:
	-go run endobit.io/builder@latest init

build::
	CGO_ENABLED=0 $(GO_BUILD) -o stack .

clean::
	rm -f stack

nuke::
	rm -rf internal/generated
