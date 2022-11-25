> 封装工具包，go-mod做包管理，依赖包版本可控。

## 组成
- **commonhelper （公共）**
  - errno：统一错误码
  - logger：统一日志 _logru_
  - tools：公共方法
  - config：统一配置读取 _viper_
- **dbhelper （数据相关）**
  - dbhelper： 数据库单例连接
  - generate： 自动生成model文件
- **servicehelper （服务相关）**
  - client：rpcx调用客户端
  - server：rpcx调用服务端
    
