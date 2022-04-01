all: build-coordinator build-compute build-storage

build-coordinator:
	@echo "building coordinator..."
	@cd cmd/coordinator && go build -o coordinator

build-compute:
	@echo "building compute-node..."
	@cd compute && RUSTFLAGS=-Awarnings cargo build --bins

build-storage:
	@echo "building storage-node..."
	@cd cmd/storage-node && go build -o storage-node

proto:
	@echo "generating proto..."
	@sh scripts/generate_proto.sh

clean-all: clean
	@rm cmd/coordinator/coordinator cmd/storage-node/storage-node -v
	@cd compute && cargo clean

clean:
	@rm cmd/storage-node/data/*/*.log -v
	@etcdctl del "slots" --prefix
	@etcdctl del "services" --prefix