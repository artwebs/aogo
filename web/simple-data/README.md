# 默认接口
* Login 登录 必传参数：appId,clientId,clientVersion
`http://localhost:8080/Login?user=admin&pwd=admin&appId=1111&clientId=111&clientVersion=1.10`
* Upload 上传 post方式 文件名为：File
`http://localhost:8080/Upload`
* Download 文件下载`http://localhost:8080/Download/xxx.xxx`

# 配置参数
## RunMode 运行模式
* dev 开发模式
* product 生产模式，该模式下必须用加密模式

## SecurityKey 数据加密的密钥
