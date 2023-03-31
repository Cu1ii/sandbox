目前沙箱依旧有 bug, 暂时定位不到 bug 的位置..

描述一下 bug 吧
## 使用 go / node 的过滤模式

运行 docker 命令如下 (请自行构建镜像)

```
docker run  -v /home/fyj/test/sandbox:/home/fyj \
--name sandbox sandbox:v0.0.1 --max_cpu_time=1000 \
--max_real_time=2000 -max_memory=67108864 --max_output_size=33554432 \
--exe_path=/home/fyj/main --error_path=/home/fyj/err.out \ 
--log_path=/home/fyj/log.out --uid=0 --gid=0 --memory_limit_check_only=0 \
--output_file=/home/fyj/out.o --input_file=/home/fyj/1.in \
--seccomp_rule_name=golang
```

 在沙箱中的运行输出如下
```
{
    "cpu_time": 0,
    "real_time": 0,
    "memory": 3178496,
    "signal": -1,
    "exit_code": 0,
    "error": 0,
    "result": 4
}
```
查看日志, 日志显示执行了指定的二进制文件

```
[DEBUG] 2023/03/31 12:35:03 log.go:62: child process
[DEBUG] 2023/03/31 12:35:03 log.go:62: Exec binary
```
后查看输出文件发现, 输出正确

------------------------------------ 分割线 ---------------------------------------------------------------
## 使用 c/cpp & general  的过滤模式

运行 docker 命令如下 (请自行构建镜像)

```
docker run  -v /home/fyj/test/sandbox:/home/fyj \
--name sandbox sandbox:v0.0.1 --max_cpu_time=1000 \
--max_real_time=2000 -max_memory=67108864 --max_output_size=33554432 \
--exe_path=/home/fyj/main --error_path=/home/fyj/err.out \ 
--log_path=/home/fyj/log.out --uid=0 --gid=0 --memory_limit_check_only=0 \
--output_file=/home/fyj/out.o --input_file=/home/fyj/1.in \
--seccomp_rule_name=c_cpp
```

在沙箱中的运行输出如下
```
{
    "cpu_time": 0,
    "real_time": 0,
    "memory": 2854912,
    "signal": 31,
    "exit_code": -1,
    "error": 0,
    "result": 4
}
```
查看日志, 日志显示并未执行指定的二进制文件, 猜测是使用 cgo 在进行资源限制时出现问题

```
fyj@iZuf6i9fhlhwd8kp71o0z4Z:~/test/sandbox$ cat log.out
fyj@iZuf6i9fhlhwd8kp71o0z4Z:~/test/sandbox$
```
后查看输出文件, 并未输出结果