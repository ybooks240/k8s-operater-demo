# Redis-Operator

## init

```
kubebuilder init --domain buleye.com --owner "James.Liu" --repo github.com/ybooks240/redisSentinel
go mod tidy
kubebuilder create api --group redis --version v1 --kind RedisSentinel --controller --resource
make 
```

## 运行

- 执行部署

```
make 
make install
kubectl apply -f 
```

- 运行控制器
```
make run 
```

- 远端运行
```
make docker-build
make docker-push
make install
```

## 测试

```
./cmd/redis --address <sentinelIP>:<sentinelPort>  -m sentinel  -p <passwd> set username james
```