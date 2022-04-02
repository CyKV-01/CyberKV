mkdir -p /tmp/cyberkv
cd ./cmd/coordinator && ./coordinator 2> /tmp/cyberkv/coordinator.log &
./compute/target/debug/compute-node -e "127.0.0.1:2379" --addr "127.0.0.1:5678" 2> /tmp/cyberkv/compute.log &
cd ./cmd/storage-node && ./storage-node -p 5796 -pp 9595 2> /tmp/cyberkv/storage-node0.log &
cd ./cmd/storage-node && ./storage-node -p 5797 2> /tmp/cyberkv/storage-node1.log &
cd ./cmd/storage-node && ./storage-node -p 5798 2> /tmp/cyberkv/storage-node2.log &
cd ./cmd/storage-node && ./storage-node -p 5799 2> /tmp/cyberkv/storage-node3.log &