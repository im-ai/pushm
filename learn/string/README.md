一、性能测试   
```
go test -bench=".*" -cpuprofile=cpu.profile .   
```


二、 进入 profile 交互模式   
```
go tool pprof string.test.exe cpu.profile   
```
top 命令
```
(pprof) top
```
web 命令 
报异常错误  找不到 dot     
直接安装 https://graphviz.gitlab.io/_pages/Download/windows/graphviz-2.38.msi
```
(pprof) web
```





