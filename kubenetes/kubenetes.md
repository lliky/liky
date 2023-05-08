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