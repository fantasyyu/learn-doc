## 查看端口
1. ss -lntp | grep ':8080'
返回LISTEN 0 4096 0.0.0.0:8080 0.0.0.0:* users:(("java",pid=1234,fd=45))
2. ps -fp 1234 查看PID进程

## 查看系统的cpu memory
1. top
2. htop
3. free -h
4. vmstat 1 5
