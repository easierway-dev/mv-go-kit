apt-get install -y redir
nohup redir --lport=8500 --caddr='consul' --cport=8500 >> /tmp/consul_port_forward.log &


for k in `ls  consul_kvs`
do
    key=`echo $k |tr "#" "/"`
    echo "[consul key]", $key
	val=`cat consul_kvs/$k`
    curl -XPUT -d "$val" "http://127.0.0.1:8500/v1/kv/$key"
done
