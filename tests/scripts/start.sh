export PATH=$PATH:/root/go/bin
pkill nchd
pkill nchcli

sleep 1
nohup nchd start &
sleep 0.1
nohup nchcli rest-server > nchcli.out 2>&1 &
