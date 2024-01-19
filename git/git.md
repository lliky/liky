# Git

## 命令

### tag

#### 列出 tag
```shell
# -l --list 可选
$ git tag

```

#### 创建 tag
```shell
$ git tag tag_name
# 创建一个带有注释的标签
$ git tag -a tag_name -m "your annotation"
```
#### 修改 tag 描述信息

```shell
$ git tag -m <your message> -e <tagname> -f
# 推送
$ git push --tags
```

#### 查看标签信息
```shell
$ git show tag_name
```
#### 删除 tag
```shell
$ git tag -d tag_name
$ git push --delete origin tag_name
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