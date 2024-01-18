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