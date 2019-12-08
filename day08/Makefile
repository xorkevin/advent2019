.PHONY: run rs run-rs build-rs

DAY=00

run:
	go run main.go

rs: build-rs run-rs

run-rs:
	./target/release/day$(DAY)

build-rs:
	cargo build --release
