prepare:
	go mod download
	go mod verify

test:
	./scripts/test.sh
