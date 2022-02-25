apt-get install -y redir 
nohup redir --lport=8500 --caddr='consul' --cport=8500 >> /tmp/consul_port_forward.log & 
