# Git

## 命令

### tag

#### 列出 tag

```shell
# -l --list 可选
$ git tag
```



#### 修改 tag 描述信息

```shell
$ git tag -m <your message> -e <tagname> -f
$ git tag push --tags -f
```

### commit

#### 回退指定 commit

```shell
# 查看日志
$ git log
# 本地回退
$ git reset --hard commit_id
# 同步远程
$ git push origin HEAD --force
```