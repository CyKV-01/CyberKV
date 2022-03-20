kill -9 $(ps -e | grep coordinator | awk '{print $1}')
kill -9 $(ps -e | grep compute-node | awk '{print $1}')
kill -9 $(ps -e | grep storage-node | awk '{print $1}')