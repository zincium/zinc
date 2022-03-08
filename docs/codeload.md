# 存储库 BLOB/ARCHIVE 缓存设计

URL 匹配，可以使用正则匹配，mux 就支持：

```
https://github.com/protocolbuffers/protobuf/raw/{branch}/{file}
```
这个可以匹配成： `/{ns}/{repo}/raw/{pathx:.+}`

由于分支名可以存在 '/'，因此我们不能简单的匹配成 `/{ns}/{repo}/raw/{branch}/...`，这样会有错误的 URL。因此这里得将分支和文件相对路径连在一起。

这样的路径如何获得对应的文件？与前后端分离的 blob 查看不一样，blob 可以传递分支信息到后端以解析出正确的引用，这里我们需要一些处理：

首先我们要明白，在存储库中，出现了 a 的分支，就不能出现 a/b 这样的分支，这是由于文件系统的限制决定的，如果在打包引用中创建了，我们就将其视为错误，基于这个前提，我们可以使用引用匹配，前缀匹配成功即可获得相应的分支和文件，如果匹配失败，我们还可以计算第一个元素是否是 SHA-1/SHA256 16 进制 ID。

```
git for-each-ref refs/ --sort=refname '--format=%(objectname) %(refname:short)/'
# de18f2411fc74d8fe2c792e638219665cd0dfe75 mainline/
```

获得文件的 BLOB id 后，我们就可以将其与项目的 ID 作为缓存 Key 存储到 codeload 缓存系统中，缓存系统添加过期时间，过期时间内可以使用，过期后被删除，命中后写入。该机制对于可以优化很多操作的时长。