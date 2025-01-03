

# Kubernetes

## 1. 认识 Kubernetes

### 1.1 为什么需要 Kubernetes

#### 1.1.1 应用部署三大阶段

1. 传统化部署：环境不隔离

2. 虚拟化部署：占用资源过多

3. 容器化部署

   > 容器：实现文件系统、网络、CPU、内存、磁盘、进程、用户空间等资源的隔离

#### 1.1.2 K8s 的特点

* 自我修复
* 弹性伸缩
* 自动部署和回滚
* 服务发现和负载均衡
* 机密和配置管理
* 存储编排
* 批处理

### 1.2 企业级容器调度平台

1. Apache Mesos
2. Docker Swarm
3. Google Kubernetes

## 2. 集群架构与组件

### 2.1相关组件

#### 2.1.1 控制面板组件

1. etcd

   键值类型的分布式数据库，提供了基于 Raft 算法实现自主的集群高可用

2. kube-apiserver

   接口服务，基于 REST 风格开放 k8s 接口的服务

3. kube-controller-manager

   管理各个类型的控制器针对 k8s 中的各种资源进行管理

   控制平面组件，负责运行控制器进程。

   控制器包括：

   1. 节点控制器：负责在节点出现故障时进行通知和相应
   2. 任务控制器：检测代表一次性任务的 Job
   3. 端点分配控制器：填充端点分片（EndpointSlice）对象
   4. 服务账号控制器：为新的命名空间创建默认的服务账号

4. cloud-controller-manager

   云控制管理器：第三方云平台提供的控制器 API 对接管理功能

5. kube-scheduler

   调度器：负责将 Pod 基于一定算法，将其调用到更合适的节点上

#### 2.1.2 节点组件

1. kubelet

   负责 Pod 的声明周期、存储、网络

2. kube-proxy

   网络代理，负责 service 的服务发现、负载均衡

3. container runtime

   容器运行时环境：docker、containerd、CRI-O

#### 2.1.3 附加组件

1. kube-dns
2. ingress controller
3. prometheus
4. dashboard
5. federation
6. elasticsearch

## 3. 核心概念

### 3.1 服务的分类

#### 3.1.1 无状态

不会对本地环境产生任何依赖

代表应用：

* nginx
* apache

优点：对客户端透明，无依赖关系，可以高效实现扩容、迁移

缺点：不能存储数据，需要额外的数据服务支撑

#### 3.1.2 有状态

会对本地环境产生依赖

代表应用：

* MySQL
* Redis

优点：可以独立存储数据，实现数据管理

缺点：集群环境下需要实现主从、数据同步、备份、水平扩容复杂

### 3.2 资源和对象

Kubernetes 中所有内容都被抽象成资源，如 Pod、Service、Node 等都是资源。“对象”就是“资源”的实例，是持久化的实体。如某个具体的 Pod、某个具体的 Node。Kubernetes 使用这些实体表示整个集群的状态。

对象的创建、修改、删除都是通过 "kubernetes API"，也就是 "Api Server" 组件提供的 API 接口，这些是 RESTful 风格的 api。

k8s 资源类别很多，kubectl 可以通过配置文件来创建这些“对象”，配置文件是描述对象”属性“的文件，配置文件格式可以是”JSON“ 或 ”YAML”。

#### 3.2.1 资源的分类

元数据级别的资源，对于资源元数据的描述，每一个资源都可以使用元空间的数据

集群级别的资源，作用于集群之上，集群下的所有资源都可以共享使用。

命名空间级别的资源，作用在命令空间之上，通常只能在该命名空间范围内使用。

##### 3.2.1.1 元数据型

* Horizontal Pod Autoscaler(HPA)

  Pod 自动扩容：可以根据 CPU 使用率或自定义指标自动对 Pod 进行扩/缩容。

* PodTemplate

  它是关于 Pod 的定义，但是被包含在其他的 k8s 对象中（如：Deployment、StatefulSet、DaemonSet 等控制器）。控制器通过 Pod Template 信息来创建 Pod。

* LimitRange

  可以对集群内 Request 和 Limits 的配置做一个全局的统一的限制，相当于批量设置了某一个范围内（某个命名空间）的 Pod 的资源使用限制。

##### 3.2.1.2 集群级

* Namespace
* Node
* ClusterNode
* ClusterRoleBinding

##### 3.2.1.3 命名空间级

###### 3.2.1.3.1 工作负载型 Pod

它是 k8s 中最小的可部署单元。一个 Pod 包含一个应用程序容器、存储资源、一个唯一的网络 IP 地址、以及一些确定容器该如何运行的选项。Pod 容器组代表了 k8s 中一个独立的应用程序运行实例。该实例可能由单个容器或几个耦合在一起的容器组成。

k8s 集群中的 Pod 存在如下两种使用途径：

* 一个 Pod 中只运行一个容器。"one-container-per-pod" 是最常见的方式。
* 一个 Pod 中运行多个需要互相协作的容器。可以将多个紧密耦合、共享资源且始终在一起运行的容器编排在同一个 Pod 中，可能的情况有：

**副本（replicas）**

一个 Pod 可以被复制多份，每一份可被称之为副本。除了描述性的信息不同（Pod 名字，uid ），其他都是一样，比如：Pod 内部的容器，容器数量，容器运行的应用。

**控制器**

* 适用无状态服务

  * ReplicationController(RC)

    帮助我们动态更新 Pod 的副本数

    不用了，用下面这个

  * **ReplicaSet(RS)**

    帮助我们动态更新 Pod 的副本数，可以通过 selector 来选择对哪些 Pod 生效

  * **Deployment**

    针对 RS 的更高层次的封装，提供了更丰富的部署相关的功能

    * 创建 Replica Set / Pod
    * 滚动升级/回滚
    * 平滑扩容和缩容
    * 暂停与恢复 Deployment

* 适用有状态服务（statefulSet）

  专门针对有状态服务进行部署的一个控制器

  * 主要特点

    * 稳定的持久化存储
    * 稳定的网络标志
    * 有序部署，有序扩展
    * 有序收缩，有序删除

  * 组成

    * Headless Service

      对于有状态服务的 DNS 管理

    * volumeClainTempalte

      用于创建持久化卷的模版

  * 注意事项

    * k8s v1.5 版本以上才支持
    * 所有 Pod 的 Volume 必须使用 PersistentVolume 或者是管理员事先创建好
    * 为了保证数据安全，删除 statefulSet 时不会被删除 Volume
    * StatefulSet 需要一个 Headless Service 来定义 DNS domain，需要在 StatefulSet 之前创建好

* 守护进程（DaemonSet）

  DaemonSet 保证在每个 Node 上都运行一个容器副本，常用来部署一些集群的日志、监控或者其他系统管理应用。典型的如下：

  * 日志收集，fluentd, logstash 
  * 监控系统 prometheus Node exporter, 
  * 系统程序 kube-proxy, kube-dns, glusterd, ceph 等

* 任务/定时任务

  * Job ：一次性任务，运行完成后 Pod 销毁，不再重新启动新容器
  * CronJob ：在 Job 基础上加上了定时功能，周期性执行

###### 3.2.1.3.2 服务发现

* Service（东西流量）

  Pod 不能直接提供外网访问，而是应该使用 service。Service 就是把 Pod 暴露出来提供服务，Service 才是真正的“服务”。

  可以说 Service 是一个应用服务的抽象，定义了 Pod 逻辑集合和访问这个 Pod 集合的策略。Service 代理 Pod 集合，对外表现为一个访问入口，访问该入口的请求将经过负载均衡，转发到后端 Pod 中的容器。

  实现 k8s 集群内部网络调用、负载均衡（四层负载）

* Ingress（南北流量）

  实现 k8s 内部服务暴露给外网访问的服务

###### 3.2.1.3.3 存储

* Volume

  数据卷，共享 Pod 中容器使用的数据。用来放持久化的数据，如：数据库数据

* CSI

  Container Storage Interface, 标准接口规范

###### 3.2.1.3.4 特殊类型配置

* ConfigMap
* Secret
  * Service Account
  * Opaque
  * kubernetes.io/dockerconfigjson
* DownwardAPI

###### 3.2.1.3.5 其他

* Role
* RoleBinding

#### 3.2.2 资源清单

### 3.3 对象规约和状态

#### 3.3.1 规约（Spec）

它描述了对象的期望状态（Desired State），希望对象所具有的特征。当创建 Kubernetes 对象时，必须提供对象的规约，用来描述该对象的期望状态，以及关于对象的一些基本信息

#### 3.3.2 状态（Status）

它表示了对象的实际状态，该属性由 k8s 自己维护，k8s 会通过一系列的控制器对对应对象进行管理，让对象尽可能的让实际状态与期望状态重合。



## 4. API 概述

### 4.1 类型

* Alpha
* Beta
* Stable

### 4.2 访问控制

* 认证
* 授权

### 4.3 废弃 API 说明

https://kubernetes.io/zh-cn/docs/reference/using-api/deprecation-guide/



## 5. 深入 Pod

### 5.1 Pod 配置文件

[Pod](./k8s资源清单.md)

### 5.2 探针

容器内应用的监测机制，根据不同的探针来判断容器应用当前的状态

#### 5.2.1 类型

#####  5.2.1.1 StartupProbe

用于判断应用程序是否已经启动

当配置了 StartupProbe 后，会先禁用其他探针，直到 startupProbe 成功后，其他探针才会继续。

作用：由于有时候不能准确预估应用一定是多长时间启动成功，因此配置另外两种方式不方便配置初始化时长来检测，而配置了 startupProbe 后，只有在应用启动成功了，才会执行另外两种探针，可以更加方便的结合使用另外两种探针。

[startupProebe](./pods/nginx-startup-probe.yaml)

##### 5.2.1.2 LivenessProbe

用于探测容器中的应用是否正常运行，如果探测失败，kubelet 会根据配置的重启策略进行重启，若没有配置，默认就认为容器启动成功，不会执行重启策略。

##### 5.2.1.3 ReadinessProbe

用于探测容器中应用是否准备好，如果准备好，就会让流量打进来

#### 5.2.2 探测方式

##### 5.2.2.1 ExecAction

在容器内部执行一个命令；如果返回值为 0，则任务容器是健康的。

##### 5.2.2.2 TCPSocketAction

通过 tcp 连接检测容器内端口是否开放，如果开放则证明该容器健康

##### 5.2.2.3 HTTPGetAction

发送 HTTP 请求到容器内的应用程序，如果接口返回状态码在 200 ～ 400 之间，则认为容器健康。

#### 5.2.3 参数配置

* initialDelaySeconds: 60   # 初始化时间
* timeoutSeconds: 2  # 超时时间
* periodSeconds: 5  # 检查时间间隔
* successThreshold: 1 # 检查 1 次成功就表示成功
* failureThreshold: 2  # 检查失败 2 次就表示失败

### 5.3 生命周期

1. 初始化阶段（可以0个或多个容器，一个一个的初始化）启动（Start 钩子函数[postStart]）
2. Pod 内的主容器（main container）
   1. 启动（Start 钩子函数[postStart])
   2. StartupProbe 启动探针
   3. readinessProbe 就绪探针/ livenessProbe 存活探针（在 Pod 之后的整个生命周期）
   4. 结束（Stop 钩子函数[preStop]）

#### 5.3.1 Pod 退出流程

删除操作

* Endpoint 删除 pod 的 IP 地址

* Pod 变成 Terminating 状态

  > 变为删除状态之后，会给 pod 一个宽限期，让 pod 去执行一些清理或销毁操作。
  >
  > 配置参数：
  >
  > terminationGracePeriodSeconds: 30
  >
  > containers:
  >
  > -- xxx

* 执行 preStop 的指令

#### 5.3.2 PreStop 的应用

* 注册中心下线
* 数据清理
* 数据销毁

## 6. 资源调度

### 6.1 Label 和 Selector

#### 6.1.1 标签（Label）

查看 labels

```shell
kubectl get pod podname --show-labels
```

如何修改？

* 配置文件

  在各类资源的 metadata.labels 中进行配置

* kubectl 

  * 临时创建 label

    ```shell
    kubectl label pod <pod-name> key=value -n namespace
    ```

  * 修改已经存在的标签

    ```shell
    kubectl label pod <pod-name> key=value2 -n namespace --overwrite
    ```

#### 6.1.2 选择器（Selector）

* 配置文件

  在各个对象的配置 spec.selector 或其他可以写 selector 的属性中编写

* kubectl

  ```shell
  # 匹配单个值，查找 app=hello 的 pod
  kubectl get pod -A -l app=hello
  # 匹配多个值
  kubectl get pod -A -l 'k8s-app in (metrics-server, kubernetes-dashboard)'
  # 多值查询
  kubectl get pod -l version!=1.2.0,type=app   # 与的关系
  # 不等值 + 语句
  kubectl get pod -l 'version!=1.2.1,type=app,k8s-app in (server, dashboard)'
  ```

### 6.2 Deployment

#### 6.2.1 功能

##### 6.2.1.1 创建

创建一个 deployment

```shell
# 创建
kubectl create deployment nginx-deploy --image=nginx

kubectl create deployment -f xxx.yaml --record
# --record 会在 annotation 中记录当前命令创建或升级了资源，后续可以查看做过哪些变动操作

# 查看部署信息
kubectl get deployment  # deploy
kubectl get replicaset # rs
```

[deploy.yaml](./deployment/nginx-deploy.yaml)

##### 6.2.1.2 滚动更新

只修改了 deployment 配置文件中的 template 中的属性后，才会触发更新操作

```shell
# 修改 nginx 版本号
kubectl set image deployment/nginx-deploy nginx=nginx:1.9.1
# edit
kubectl edit deployment deployment-name
# 查看滚动更新的过程
kubectl rollout status deployment deployment-name
```

滚动更新并行

> 会按照最后一次的修改进行滚动更新  

##### 6.2.1.3 回滚

在默认情况下，k8s 会在系统中保存前两次的 Deployment 的 rollout 历史记录，以便随时回退。（可以修改 revision history limit 来更改保存的 revision 数）

```shell
# 查看revision 列表
kubectl rollout history deploy deploy-name 
# 查看详细信息
kubectl rollout history deploy deploy-name --revision=id
# 回退
kubectl rollout undo deploy deploy-name --to-revision=id
	
```

可以通过设置 spec.revisionHistoryLimit 来指定保留多少 revision ，如果设置为 0，则不允许 deployment 回退。

##### 6.2.1.4 扩容缩容

通过 kubectl scale 命令可以进行自动扩缩容，以及通过 kubectl edit 编辑 replicas 也可以实现。

扩缩容只是创建副本数，没有更新 pod template 因此不会创建新的 rs。

```shell
kubectl scale --help  #命令查看
```

##### 6.2.1.5 暂停与恢复

由于每次对 pod template 中的信息发生修改后，都会触发更新 deployment 操作，那么如果频繁修改信息，就会产生多次更新，而实际上只需要执行最后一次更新即可，当出现此类情况，我们就可以暂停 deploy 的 rollout

 ```shell
 # 暂停
 kubectl rollout pause deploy deploy-name
 # 恢复
 kubectl rollout resume deploy deploy-name
 ```



#### 6.2.1配置文件

[deploy.yaml](./deployment/nginx-deploy.yaml)

### 6.3 StatefulSet

#### 6.3.1 功能

##### 6.3.1.1 创建

[statefulSet](./statefulSet/web.yaml)

##### 6.3.1.2 扩容缩容

```shell
kubectl scale statefulset web --replicas=5
kubectl pathc statefulset web -p '{"spec":{"replicas": 3}}' 
# 删除有顺序性

```

##### 6.3.1.3 镜像更新

暂时不支持直接更新 image，需要 patch 来间接实现

```
kubectl patch statefulset web --type='json'='["op": "replace", "path": "/spec/tempalte/spec/containers/0/image", "value":"nginx:1.9.1"]'	
```

* RollingUpdate

  滚动更新策略，同样是修改 pod template 属性后会触发更新，但是由于 pod 是有序的，在 statefulset 中更新时是基于 pod 的顺序**倒序更新**的

  **灰度发布**(金丝雀发布)

  利用滚动更新中的 partition 属性，可以实现简易的灰度发布的效果。

  例如：有 5 个 pod ，如果当前 partition 设置为 3 ，那么此时滚动更新时，只会更新那些序号 >=3 的 pod。

  利用该机制，可以可以通过控制 partition 的值，来决定只更新其中一部分 pod ，确认没问题后，在增大更新 pod 的数量，最终实现全部 pod 更新。

* OnDelete

  删除镜像之后才更新

##### 6.3.1.4 删除

* 级联删除：删除 sts 时，会同时删除 pods

* 非级联删除：删除 sts 时不会删除 pods

  ```shell
  kubectl delete sts web --cascade=false # --cascade=orphan
  ```

#### 6.3.2 配置文件

### 6.4 DaemonSet

#### 6.4.1 指定 node 节点

DaemonSet 会忽略 Node 的 unschedulable 状态，有两种方式来指定 Pod 只运行在指定的 Node 节点上

* nodeSelector：只调度到匹配指定 label 的 Node 上
* nodeAffinity：功能丰富的 Node 选择器，比如支持集合操作
* podAffinity：调度到满足条件的 Pod 所在的 Node 上

#### 6.4.2 滚动更新

不建议使用 RollingUpdate, 建议使用 OnDelete 模式，这样避免频繁更新 ds

### 6.5 HPA 自动扩/缩容

通过观察 pod 的 cpu 、内存使用率或自定义 metrics 指标进行自动的扩容或缩容 pod 的数量

通常用于 Deployment ，不适用于无法扩缩容的对象，如 DaemeonSet

控制管理器每隔 30s 查询 metrics 的资源使用情况

#### 6.5.1 cpu、内存指标监控

> 前提：该对象必须配置 resources.request.cpu 或 resources.request.memory 才可以。配置当 cpu/memory 达到上述配置的百分比后进行扩容或缩容

创建一个 HPA:

1. 准备好一个有资源限制的 deployment

2. 执行命令

   ```shell
   # --cpu-percent cpu 使用率占20 就扩容，最小两个，最大 5 个
   kubectl autoscale deploy deploy_name --cpu-percent=20 --min=2 --max=5
   ```

3. 通过 kubectl get hpa 获取 HPA 信息

#### 6.5.2 自定义 metrics



## 7. 服务发布

### 7.1 服务发现

#### 7.1.1 Service

Pod 不能直接提供给外网访问，而是应该使用 service。Service 就是把 Pod 暴露出来提供服务，Service 才是真正的“服务”。

Service 是一个应用服务的抽象，定义了 Pod 逻辑集合和访问这个 Pod 集合的策略。Service 代理 Pod 集合，对外表现为一个访问入口，访问该入口的请求将经过负载均衡，转发到后端 Pod 中的容器。

##### 7.1.1.1 Service 的定义

* 命令操作

  ```shell
  # 创建 service
  kubectl create -f file.yaml
  
  # 查看 service 信息，通过 service 的 cluster ip 进行访问
  kubectl get svc
  
  # 查看 pod 信息，通过 pod 的 ip 进行访问
  kubectl get pod -o wide
  
  # 创建其他 pod 通过 service name 进行访问
  kubectl exec -it busybox -- sh
  curl http://nginx-svc
  
  # 默认在当前 namespace 中访问，如果需要跨 namespace 访问 pod，则在 service name 后面加上 <.namespace> 即可
  curl http://gninx-svc.default
  ```

* Endpoint

  ```
  kubectl get endpoints
  ```

##### 7.1.1.2 代理 k8s 外部服务

实现方式：

1. 编写 service 配置文件时，不指定 selector 属性
2. 自己创建 endpoint

需要访问外部服务的情形：

* 各种环境访问名称统一
* 访问 k8s 集群外的其他服务
* 项目迁移

##### 7.1.1.3 反向代理外部域名

```yaml
apiVersion: v1
kind: Service
metadata:
	labels:
		app: external-domain
	name: external-domain
spec:
	type: ExternalName
	externalName: www.baidu.com
```



##### 7.1.1.4 常用类型

* ClusterIP

  只能在集群内部使用，不配置类型，就默认 ClusterIP

* ExternalName

  返回定的 CNAME 别名，可以配置为域名

* NodePort

  会在所有安装了 kube-proxy 的节点都绑定一个端口，此端口可以代理至对应的 Pod，集群外部可以使用任意节点 ip + NodePort 的端口号访问到集群中对应 Pod 中的服务

  当类型设置为 NodePort 后，可以在 ports 配置中增加 nodePort 配置指定端口，如果不指定会随机指定端口

  端口范围：30000 ～ 32767

* LoadBalancer

  使用云服务商提供的负载均衡服务



### 7.2 ingress



## 8 配置与存储

### 8.1 配置与管理

#### 8.1.1 ConfigMap

* 创建

  ```shell
  $ kubectl create configmap -h
  # 根据文件夹创建，会把文件夹里面的文件都加载到 cm
  $ kubectl create configmap my-config --from-file=path/to/bar
  
  # 指定文件去创建，常用，key 把原文件重命名
  kubectl create configmap my-config --from-file=key1=/path/to/bar/file1.txt --from-file=key2=/path/to/bar/file2.txt
  
  # 自己写 key=value 键值对
  kubectl create configmap my-config --from-literal=key1=config1 --from-literal=key2=config2
  ```

  使用 -h 查看示例，构建 ConfigMap 对象

* 使用 ConfigMap

#### 8.1.2 加密数据配置 Secret

用于存储配置信息，主要用于存储敏感信息、需要加密的信息，Secret 可以提供数据加密、解密功能。

在创建 Secret 时，需要注意还有特殊字符需要使用转义字符或单引号描述。

```shell
$ kubectl create secret -h
$ kubectl create secret generic demo-secret --from-literal=username=admin --from-literal=passwd='123'
```

#### 8.1.3 SubPath 的使用

在使用 ConfigMap 或 Secret 挂载到目录的时候，会将容器中源目录给覆盖掉，但是我们只想覆盖目录中一个文件，这时就需要使用 SubPath。

配置方式：

1. 定义 volumes 时需要增加 items 属性，配置 key 和 path，且 path 不能以 '/' 开始
2. 在容器内的 volumeMounts 中增加 subPath 属性，该值与 volumes 中 items.path 的值相同

#### 8.1.4 配置数据的热更新

1. 通过 edit 命令直接修改 configmap

   通常将项目的配置文件作为 configMap，然后挂载到 pod ，那么如果更新 configMap 中的配置，会不会更新到 pod 中呢？

   几种情况：

   * 默认方式：会更新，更新周期是更新时间 + 缓存时间
   * subPath：不会更新
   * 变量形式：如果 pod 中的一个变量是从 configMap 或 secret 中得到，不会更新

2. 通过 replace 替换

   ```
   kubectl create cm --from-file=nginx.conf --dry-run -o yaml | kubectl replace -f-
   ```

   

#### 8.1.5 不可变的 Secret 和 ConfigMap

在配置 configmap 时可以设置 immutable: true 来禁止修改

### 8.2 持久化存储

#### 8.2.1 Volumes

##### 8.2.1.1 HostPath

将节点上的文件或目录挂载到 Pod 上，此时该目录会变成持久化存储目录，即使 Pod 被删除后重启，也可以重新加载到该目录，文件不会丢失。

##### 8.2.1.2 EmptyDir

主要用于一个 Pod 中不同的 Container 共享数据使用的，由于只是在 Pod 内部使用，因此与其他 volume 比较大的区别是，当 Pod 如果被删除了，那么 emptyDir 也会被删除。

存储介质可以是任意类型，如 SSD、磁盘或网络存储。

#### 8.2.2 NFS 挂载

NFS 卷能将 NFS(网络文件系统) 挂载到你的 Pod 中。nfs 卷的内容在删除 Pod 时会被保存，卷只是被卸载。意味着 nfs 可以被预先填充数据，并且这些数据可以在 Pod 之间共享。

需要安装 NFS 服务

#### 8.2.3 PV 与 PVC

##### 8.2.3.1 生命周期

##### 8.2.3.2 PV

##### 8.2.3.3 PVC

##### 8.2.3.4 StorageClass



## 9 高级调度

### 9.1 CronJob 计划任务

在 k8s 中周期性运行计划任务。

#### 9.1.1 表达式

```
------ 分钟（0 - 59）
| ------ 小时（0 - 23）
| | ------ 月的某天（1 - 31）
| | |	------ 月份（ 1 - 12）
| | | | ------ 周的某天（0 - 6）
| | | | |
| | | | |
* * * * * 
```

### 9.2 初始化容器 InitContainer

### 9.3 污点和容忍

#### 9.3.1 污点（taint）

* NoSchedule：如果不能容忍该节点，那么 pod 就无法调度到该节点上
* NoExecute：不能执行
  * 如果 Pod 不能忍受污点，Pod 会马上被驱逐
  * 如果 Pod 能忍受这类污点，但是没有指定 tolerationSeconds ，则 Pod 还会一直在该节点上运行
  * 如果 Pod 能忍受这类污点，且指定 tolerationSeconds ，则 Pod 还在该节点上运行指定的时间长度

```shell
# 添加污点
kubectl taint node node-name key=value:NoSchedule
# 移除污点
kubectl taint node node-name key=value:NoSchedule-
# 查看污点
kubectl describe node node-name
```

#### 9.3.2 容忍度（toleration）

标注在 pod 上的，当 pod 被调度时，如果没有配置容忍，则该 pod 不会被调度到有污点的节点上，只有该 pod 上标注了满足某个节点的所有污点，则会被调度到这些节点

```yaml
# pod 的 spec 下面配置容忍
toleration:
- key: "污点的 key "
	value: "污点的 value "
	offect: "NoSchedule" # 污点产生的影响
	operator: "Equal/Exists"
```

* Equal：意味着必须与污点值做匹配，key/value 都必须相同，才表示能够容忍该污点
* Exists：容忍与污点的比较只比较 key，不关心 value 的值

### 9.4 亲和力（Affinity）

#### 9.4.1  NodeAffinity

#### 9.4.2 PodAffinity

#### 9.4.3 PodAntiAffinty
