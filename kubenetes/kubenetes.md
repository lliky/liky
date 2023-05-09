## Service





## 命令
```
# 添加标签
kubectl label 资源  资源名  key=value
kubectl label pod pod-name app=web 
# 删除标签
kubectl label 资源 资源名 key-
kubectl label pod pod-name app-
# 查看版本
kubectl rollout history deployment app
# 回滚
kubectl rollout undo deployment app --to-revision=1
```
# Docker
## Dockerfile 指令
```
FROM 继承基础镜像
MAINTAINER 镜像制作作者信息
RUN 用来执行 shell 命令
EXPOSE 暴露端口
CMD 启动容器默认执行的命令
ENTRYPOINT 启动容器真正执行的命令
# CMD 和 ENTRYPOINT 必须要有一个
# CMD 可以被覆盖，如果有 ENTEYPOINT 的话，CMD 就是 ENTRYPOINT 的参数
# 可以覆盖 CMD 命令  docker run -it imageName cover, cover 就能覆盖CMD, 相当于 CMD= cover ，可作为 ENTRYPOINT 参数
VOLUME 创建挂载点
ENV 配置环境变量
ADD 复制文件到容器
COPY 复制文件到容器
## ADD 会解压文件  COPY 不会解压
WORKDIR 设置容器的工作目录
USER 容器使用的用户
```

### 制作小镜像
使用多阶段构建，编译操作和生成最终镜像的操作  FROM, FROM

镜像 scratch 是一个空镜像

# Kubernetes

## Master 节点
整个集群的控制中枢
* Kube-APIServer : 集群的控制中枢，各个模块之间信息交互都需要经过 Kube-APIServer，同时它也是集群管理、资源配置、整个集群安全配置的入口。  
* Controller-Manager: 集群的状态管理器，保证 Pod 或其他资源达到期望值，也是需要和 APIServer 进行通信，在需要的时候创建、更新或删除它所管理的资源。  
* Scheduler: 集群的调度中心，它会根据指定的一系列条件，选择一个或一批最佳的节点，然后部署 Pod。  
* Etcd: 键值数据库，保存一些集群的信息。建议部署三个以上奇数节点。

## Node 节点
* Kubelet: 负责监听节点上 Pod 的状态，同时负责上报节点和节点上面 Pod 的状态，负责与 Master 节点通信，并管理节点上面的 Pod。
* Kube-Proxy: 负责 Pod 之间的通信和负载均衡，将指定的流量分发到后端正确的机器上。