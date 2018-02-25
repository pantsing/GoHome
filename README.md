GoHome
======

GoHome helps connect your home computers/servers to internet via DNSPOD Dynamic DNS.

Installation
------
<pre><code>
git clone https://github.com/pantsing/GoHome.git

cd gohome

# 编译生成待安装文件 gh.tar.gz

sh dist.sh

# 解压部署

tar -C /usr/local/ -xzf gh.tar.gz

# 根据您的DNSPOD配置及要动态更新IP地址的域名修改配置文件

vim /usr/local/gohome/gohome.toml

#启动服务

cp /usr/local/gohome/*.service /usr/lib/systemd/system/

systemctl enable ghclient.service

systemctl start ghclient.service
</code></pre>
