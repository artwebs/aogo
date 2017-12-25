### 日志
```
1、app.conf中设置日志等级 Debug = 0 #0为生产环境 1为测试环境
2、临时调试，在使用方法中加入
  logger.SetLevel(logger.LoggerDebug)
  defer logger.SetLevel(web.LoggerLevelDefault)
3、调试日志：
  logger.Debug("debug")
  logger.Info("Info")
  logger.Notice("Notice")
  logger.Warn("Warn")
  logger.Error("Error")
  logger.Fatal("Fatal")
```
### 数据表操作
```
1、新建操作对象：table := db.TableNew("表名", "前缀") #默认前缀是bas_
2、查询：table.Where("belongs_month=?", val).Select()
```
#### 默认加解密
```
#sn 加密序号
#val加密值
#desObj加密值如security.NewSecurityDES()  
#secrets 密钥 map[string]security.Secret{"01": security.Secret{Key: "xxxxxxxx", Iv: "xxxxxxxx"}}
1、加密：security.Encrypt(sn, val, desObj, secrets)
2、解密：security.Decrypt(val, desObj, secrets)
```
