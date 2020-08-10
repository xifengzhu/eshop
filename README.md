## Eshop
### Project Structure
```
├── Dockerfile
├── README.md
├── conf
│   ├── permissions.json
│   └── rbac_model.conf
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── helpers
│   ├── e
│   ├── export
│   └── utils
├── initializers
│   ├── captcha.go
│   ├── permission.go
│   ├── qiniu.go
│   ├── redis_pool.go
│   ├── setting
│   ├── wechat
│   └── wechat.go
├── log
├── main.go
├── middlewares
│   ├── ip_filter
│   ├── jwt
│   ├── logger
│   └── role
├── models
│   ├── address.go
│   ├── admin_user.go
│   ├── base_models.go
│   ├── ...
├── public
│   ├── export
│   └── qrcode
├── routers
│   ├── admin_api/
│   ├── app_api/
│   ├── helpers/
│   ├── validators/
│   └── router.go
├── seeds
│   └── data
├── tmp
│   ├── air_errors.log
│   └── main
├── uploads
├── workers
│   ├── export_job.go
│   ├── order_job.go
│   ├── refresh_wechat_token_job.go
│   └── worker.go
├── go.mod
└── go.sum
```

### Start project
```
air -c .air.conf
```

### Start go worker webui
```
workwebui -redis="redis:6379" -ns="eshop" -listen=":5040"
```

### API Doc
#### url
visit: http://127.0.0.1:8000/swagger/index.html

#### Update swagger document
```
swag init
```