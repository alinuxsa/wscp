## 构建

```
go build .
```

## 服务端监听地址

```
wscp s 0.0.0.0:5467
```

## 客户端连接服务端并上传dist.zip

```
wscp c your_domain_here -f dist.zip
```

## 使用Nginx反向代理

```
map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

upstream websocket {
    server wscp_ip_addr:5467; 
}

server {
     server_name your_domain_here;
     listen 80;
     location /wscp {
         proxy_pass http://websocket;
         proxy_read_timeout 300s;
         proxy_send_timeout 300s;
         
         proxy_set_header Host $host;
         proxy_set_header X-Real-IP $remote_addr;
         proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
         
         proxy_http_version 1.1;
         proxy_set_header Upgrade $http_upgrade;
         proxy_set_header Connection $connection_upgrade;
     }
}
```