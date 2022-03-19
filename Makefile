all: build-coordinator build-compute

build-coordinator:
	@echo "building coordinator..."
	@cd cmd && go build -o coordinator

build-compute:
	@echo "building compute-node..."
	@cd compute && cargo build --bins

clean:
	@rm cmd/coordinator
	@cd compute && cargo clean