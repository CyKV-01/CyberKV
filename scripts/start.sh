cd ./cmd/coordinator && ./coordinator 2> /tmp/coordinator.log &
./compute/target/debug/compute-node -e "127.0.0.1:2379" --addr "127.0.0.1:5678" 2> /tmp/compute.log &
cd ./cmd/storage-node && ./storage-node 2> /tmp/storage-node.log &