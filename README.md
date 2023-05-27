
1.允许多端登录
2.多端踢出下线
3.消息数据解压缩
4.唯一的operationID,设置追踪
5.在线消息回调

问题:
1.为什么只读,不需要监听写
2.已读的回执怎么做

nginx配置:
server {
listen 50001;
server_name localhost;
gzip on;
gzip_min_length 1k;
gzip_buffers 4 16k;
gzip_comp_level 2;
gzip_types text/plain application/javascript application/x-javascript text/css application/xml text/javascript application/x-httpd-php image/jpeg image/gif image/png;
gzip_vary off;
gzip_disable "MSIE [1-6]\.";

        location / {
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection "Upgrade";
                proxy_set_header X-real-ip $remote_addr;
                proxy_set_header X-Forwarded-For $remote_addr;
                proxy_pass http://localhost:10001;
            }
}

