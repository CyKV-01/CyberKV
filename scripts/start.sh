mkdir -p /tmp/cyberkv
mkdir -p ./tmp/coordinator
mkdir -p ./tmp/compute1
mkdir -p ./tmp/storage1
mkdir -p ./tmp/storage2
mkdir -p ./tmp/storage3
cd ./tmp/compute1 && ../../compute/target/release/compute-node -e "127.0.0.1:2379" --addr "127.0.0.1:5678" 2> /tmp/cyberkv/compute.log &
cd ./tmp/storage1 && ../../cmd/storage-node/storage-node -p 5796 -pp 9595 2> /tmp/cyberkv/storage-node0.log &
sleep 0.5
# cd ./tmp/storage2 && ../../cmd/storage-node/storage-node -p 5797 2> /tmp/cyberkv/storage-node1.log &
# sleep 0.5
# cd ./tmp/storage3 && ../../cmd/storage-node/storage-node -p 5798 2> /tmp/cyberkv/storage-node2.log &
# sleep 0.5
cd ./tmp/coordinator && ../../cmd/coordinator/coordinator 2> /tmp/cyberkv/coordinator.log &
