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
    ```
    查看 Kube-proxy 工作模式： curl 127.0.0.1:10249/proxyMode
    Ipvs: 监听 Master 节点增加和删除 service 以及 endpoint 的消息，调用 Netlink 接口创建相应的 IPVS 规则，通过 IPVS 规则，将流量转发至相应的 Pod 上。
    Iptables: 监听 Master 节点增加和删除 service 以及 endpoint 的消息，对于每一个 service ，他都会创建一个 iptables 规则，将 service 的 clusterIP 代理到后端对应的 Pod。
    ```
* CoreDNS: 用于 Kubernetes 集群内部 Service 的解析，可以让 Pod 把 Service 名称解析成 IP 地址，然后通过 Service 的 IP 地址进行连接到对应的应用上。

## Pod

### 什么是 Pod
Pod 是 Kubernetes 中最小的单元，它由一组、一个或多个容器组成，每个 Pod 还包含了一个 Pause 容器，Pause 容器是 Pod 的父容器，主要负责僵尸进程的回收管理，通过 Pause 容器可以使同一个 Pod 里面的多个容器共享存储、网络、PID、IPC 等。

### 定义一个 Pod
```yaml
apiVerson: v1 # 必选， API 的版本号
kind: Pod # 必选，类型 Pod
metadata: # 必选，元数据
  name: # 必选，符合 RFC 1035规范的 Pod 名称
  namesapce: default # 可选，Pod 所在的名称空间，不指定默认为 default, 可以使用 -n 指定namespace
  labels: # 可选，标签选择器，一般用于过滤和区分 Pod
    app: nginx
    role: frontend  # 可以写多个
  annotations: # 可选，注释列表，可以写多个
    app: nginx
spec: # 必选，用于定义容器的详细信息
  initContainers: # 初始化容器，在容器启动之前执行的一些初始化操作
  - command:
    - sh
    - -c
    - echo "I am InitContainer for init some configuration"
    image: busybox
    imagePullPolicy: IfNotPresent
    name: init-container
  containers: # 必选，容器列表
  - name: nginx # 必选，符合 RFC 1035规范的容器名称
    image: nginx:latest # 必选，容器所用的镜像的地址
    imagePullPolicy: Always # 可选，镜像拉取策略, IfNotPresent, Always, Never
    command: # 可选，容器启动执行的命令
    - nginx
    - -g
    - "daemon off;"
    workingDir: /usr/share/nginx/html # 可选，容器的工作目录
    volumeMounts: # 可选，存储卷配置，可以配置多个
    - name: webroot # 存储卷名称
      mountPath: /usr/share/nginx/html # 挂载目录
      readOnly: true  # 只读
    ports: # 可选，容器需要暴露的端口号列表
    - name: http # 端口名称
      containerPort: 80 # 端口号
      protocol: TCP # 端口协议，默认 TCP
    env:  # 可选，环境变量配置列表
    - name: TZ # 变量名
      value: Asia/Shanghai # 变量的值
    - name: LANG
      value: en_US.utf8
    resources: # 可选，资源限制和资源请求限制
      limits: # 最大限制设置
        cpu: 1000m
        memory: 1024Mi
      requests: # 启动所需的资源
        cpu: 100m
        memory: 512Mi
#    startupProbe: #  可选，检测容器内进行是否完成启动。注意三种检查方式同时只能使用一种
#      httpGet: # httpGet 检测方式，生产环境建议使用 httpGet 实现接口级健康检查，健康检查由应用程序提供
#        path: /api/successStart  # 检查路径
#        port: 80
    readinessProbe: # 可选，健康检查。注意三种检查方式同时只能使用一种
      httpGet: # httpGet 检测方式，生产环境建议使用 httpGet 实现接口级健康检查，健康检查由应用程序提供
        path: /  # 检查路径
        port: 80
    livenessProbe: # 可选，健康检查
#      exec: # 执行容器命令检测方式
#        command:
#        - cat
#        - /health
#      httpGet: # httpGet 检测方式
#        path: /_health # 检查路径
#        prot: 8080
#        httpHeaders: # 检查的请求头
#        - name: end-user
#          value: Jason
      tcpSocket: # 端口检测方式
        port: 80
      initialDelaySeconds: 60 # 初始化时间
      timeoutSeconds: 2 # 超时时间
      periodSeconds: 5 # 间隔时间
      successThreshold: 1 # 检测成功 1 次表示就绪
      failureTHreshold: 2 # 检测失败 2 次表示未就绪
    lifecycle:
      postStart: # 容器创建完成后执行的命令，可以是 exec, httpGet, TCPSocket
        exec:
          command:
          - sh
          - -c
          - 'mkdir /data/'
      preStop:
        httpGet:
          path: /
          port: 80
  restartPolicy: Always # 可选，默认为 Always, 容器故障或者没有启动成功，那就自动该容器；Onfailure: 容器以不为 0 的状态终止，自动重启该容器；Never, 无论什么状态，都不重启
  nodeSelector: # 可选，指定 Node 节点
    region: subnet7
  imagePullSecrets: # 可选，拉取镜像使用的 secret, 可以配置多个
  - name: default-dockercfg
  hostNetwork: false # 可选，是否为主机模式，如是会占用主机端口
  volumes: # 共享存储卷
  - name: webroot # 名称，与上对应
    emptyDir: {} # 挂载目录
#      hostPath: # 挂载本机目录
#        path: /etc/hosts
```