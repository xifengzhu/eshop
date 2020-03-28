### start project
air -c .air.conf

### start go worker webui
workwebui -redis="redis:6379" -ns="eshop" -listen=":5040"

### update swagger doc
swag init

### Feature
[Done]1. 鉴权jwt
[Done]2. CRUD
[Done]3. pagination
[Done]4. 嵌套保存方案
[Done]5. 集成swagger
[Done]6. background job&schedule job
[Done]7. air 监听代码变化自动编译

# TODO
<!--  优惠券 -->
<!--  营销模块 -->
