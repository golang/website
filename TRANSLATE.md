# 翻译流程
## 基础工作
1. fork [`https://github.com/studygolang/go`](https://github.com/studygolang/go)的代码到自己的仓库。
2. clone go项目 到本地之后，切换到`release-branch.go1.15` 分支
3. clone 本项目(website) 到本地
```bash
git clone https://github.com.cnpmjs.org/studygolang/website.git
git clone https://github.com.cnpmjs.org/your-name/go.git
cd go
git checkout release-branch.go1.15
```

## 提翻译issue
1. 在[`https://github.com/studygolang/go`](https://github.com/studygolang/go)中 `doc`目录下找想要翻译的模块或者文件
  同时去issue中查看是否有人已经在翻译了，如果没有就可以翻译.
2. 在本项目中(website) 发起一个issue,名字类似于:`翻译：标准库image`等表明自己需要翻译的文档。

## 开始翻译
1. 在go `release-branch.go1.15`的基础上 检出一个分支，在该分支下做文档翻译 `git checkout -b your-branch`
2. 翻译规范请参照[`https://github.com/studygolang/GCTT/blob/master/chinese-copywriting-guidlines.md`](https://github.com/studygolang/GCTT/blob/master/chinese-copywriting-guidlines.md)
2. 可以在 webiste目录下,启动本地server，在`localhost:6060`下查看自己的文档翻译效果 `go run ./cmd/golangorg -goroot=../go`
3. push到自己的远端

## 提pr
1. 在studygolang/go项目中，发起pull request,等待合入。

## 建议
1. 最好一个模块单独一个issue和pr
2. 建议先启动server看看需要改动哪些，_test.go,未导出变量等等均可不用翻译。
3. 如果发现还有需要改动的地方，也可以先转为draft,等修改完成之后再合入。
3. 提pr之后如果有问题，可以直接在原pr上改，最好别重新提pr。
