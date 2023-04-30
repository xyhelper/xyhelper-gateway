# xyhelper-gateway

xyhelper 接口中转网关

## 使用

创建 docker-compose.yml

```yaml
version: '3'
services:
  xyhelper-gateway:
    image: xyhelper/xyhelper-gateway:latest
    container_name: xyhelper-gateway
    restart: always
    ports:
      - 8080:8080
    environment:
      BASEURI: "https://freechat2.xyhelper.cn"
      TOKENS: "token1,token2,token3"

```

运行

```bash
docker-compose up -d
```




## 环境变量

| 环境变量 | 说明 | 默认值 |
| -------- | ---- | ------ |
| PORT     | 端口 | 8199   |
| BASEURI  | 基础地址 | https://freechat.xyhelper.cn |
| TOKENS   | 令牌列表 | xyhelper-gateway-default-token |
