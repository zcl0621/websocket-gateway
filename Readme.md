# websocket gateway

a golang websocket gateway

## Features
* horizontal scaling
* single client push
* group push
* client can get miss message
* client can get miss group message
* no ack
* heartbeat text message:ping/Ping/pong/Pong control message: ping/Pong

## Used Packages
* gin github.com/gin-gonic/gin
* gorilla/websocket github.com/gorilla/websocket
* go-redis github.com/go-redis/redis
* cors github.com/gin-contrib/cors
* xid github.com/rs/xid
* concurrent-map https://github.com/orcaman/concurrent-map

## Deployment Example
![img.png](img.png)

> nginx/ingress must set proxy-read-timeout and proxy-send-timeout to up to 300s

## Process
![img_2.png](img_2.png)