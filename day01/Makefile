.PHONY: run rs run-rs build-rs

DAY=01

run:
	go run main.go

rs: build-rs run-rs

run-rs:
	./target/release/day$(DAY)

build-rs:
	cargo build --release
