# 日志
```
1、app.conf中设置日志等级 Debug = 0 #0为生产环境 1为测试环境
2、临时调试，在使用方法中加入
  logger.SetLevel(logger.LoggerDebug)
  defer logger.SetLevel(web.LoggerLevelDefault)
```
# aogo
# aogo
