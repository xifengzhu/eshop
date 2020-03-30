## Eshop
### Project Structure
```
├── README.md
├── conf
│   └── rbac_model.conf
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── helpers
│   ├── e
│   ├── setting
│   ├── utils
│   └── wechat
├── main.go
├── middleware
│   ├── jwt
│   └── role
├── models
│   ├── base_models.go
│   ├── casbin_rule.go
│   └── crud.go
├── routers
│   ├── admin_api
│   ├── api_helpers
│   ├── app_api
│   └── router.go
├── seeds
│   └── data
└── workers
    ├── refresh_wechat_token_job.go
    ├── send_email_job.go
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

#### Update document
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
