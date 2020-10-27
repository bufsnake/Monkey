## Monkey

基于 nmap + masscan + socket port scan

## 使用场景

建议内网，路由跳为0的时候使用

`基于线程为255的情况下`

### socket扫描

```bash
理论最大耗时不会超过23分钟
在未设置网关的情况下使用
```

### masscan扫描

```
根据网络情况，耗时不会超过5分钟
有网关的情况下使用
```

## Usage

```bash
Usage of ./Monkey_macos_amd64:
  -f string
    	从文件中获取IP
  -i string
    	指定IP
  -m	指定是否使用masscan进行端口扫描
  -p int
    	指定端口扫描线程,默认50 (default 5)
  -t int
    	指定线程,默认50 (default 50)
  -v string
    	指定-sV详细程度0-9 (default "0")
```
