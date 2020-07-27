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
├── eshop
├── export
├── gin-bin
├── go.mod
├── go.sum
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
├── middleware
│   ├── ip_filter
│   ├── jwt
│   ├── logger
│   └── role
├── models
│   ├── address.go
│   ├── admin_user.go
│   ├── base_models.go
│   ├── ...
├── routers
│   ├── admin_api/
│   ├── api_helpers/
│   ├── app_api/
│   └── router.go
├── seeds
│   └── data
├── tmp
│   ├── air_errors.log
│   └── main
├── uploads
└── workers
    ├── export_job.go
    ├── order_job.go
    ├── refresh_wechat_token_job.go
    └── worker.go
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

### Points
- [x] 鉴权jwt
- [x] CRUD
- [x] pagination
- [x] 嵌套保存方案
- [x] 集成swagger API document
- [x] background job&schedule job
- [x] air 监听代码变化自动编译
- [ ] 优惠券
- [ ] 营销模块