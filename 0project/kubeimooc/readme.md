

# web框架

安装 `gin 1.10.0`
```shell
go get -u github.com/gin-gonic/gin
```

gin github 主页
<https://github.com/gin-gonic/gin>



# 配置参数分离
<https://github.com/spf13/viper>

```shell
go get github.com/spf13/viper
```

# k8s客户端库


```shell
go get k8s.io/client-go@latest
```
<https://github.com/kubernetes/client-go>


## 项目接口开发

### pod管理接口开发


- 命名空间列表接口
  - 命名空间创建 - 不做
- Pod创建
- Pod查看（详情、列表）
- Pod编辑（更新/升级）
- Pod删除