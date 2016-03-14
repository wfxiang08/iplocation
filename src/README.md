# 糖尿病400电话服务

## 1. 环境搭建:(采用gvp + gpm来实现类似Python Virtualenv的项目环境管理)
* [前端构建](https://git.chunyu.me/health/react-yaoguanjia/blob/master/README.md)
* `安装工具`
	* brew install gvp
	* brew install gpm
	* brew tap pote/gpm_plugins
	* brew install gpm-bootstrap

* 工作目录: 约定使用 `~/goprojects/`
	* mkdir diabetics400 && cd diabetics400
	* mkdir src && mkdir bin && mkdir pkg
* `下载代码`
	* cd ~/goprojects/
	* git clone git@git.chunyu.me:health/diabetics400.git
	* cd diabetics400
	* source gvp
	* cd src && gmp install (按照Godeps中的配置来布置环境)
* 补充学习: `创建项目`(项目是如何创建的)
	* beego框架:
		* 参考文档: http://beego.me/quickstart
		* go get github.com/astaxie/beego
		* go get github.com/beego/bee
		* bee new diabetics400 && mv diabetics400/* src/
    * 准备初始的Godeps
	    * gpm bootstrap
    * 启用项目私有的$GOPATH
	    * source gvp (在目录: ~/goprojects/diabetics400目录下使用)
	    * echo $GOPATH 

* IDE(`Pycharm + Go plugin`)
	* https://github.com/go-lang-plugin-org/go-lang-idea-plugin
	* 在最新的Pycharm中安装:
		* https://plugins.jetbrains.com/plugin/5047 (版本: 0.9.748)
    * 配置:
	    * 在Pycharm的 Preferences 中选择:
		    * Languages & Frameworks
		    * 选择go
		    * 设置Go SDK(设置GORoot)
		    * 选择Go Libraries, 选择Global Libraries(忽略), 使用Project libraries)
	    * 保证Pycharm的配置和source gvp的设置是一致的(在命令行可以编译，在IDE也可以编译)
		    * 这就是为什么在最后时刻我们放弃: LiteIDE

## 开发流程:
1.本地环境搭建
2.本地开发，编译，测试, 成功之后提交git
```bash
	git checkout -b feature/new_feature
	git add xxx
	git commit -a -m"添加xxx修改"
	git push origin feature/new_feature
```
3.登陆dk1
```bash
    cd ~/goprojects/diabetics400
    source start_env.sh
    ./update.sh

    ./build.sh
    # 如果出现 cannot find package，则可能是: src/Godeps 中添加了新的包依赖，或者代码中import了新的包，但是没有在Godeps中添加信息
    gpm install
	# 包更新完毕之后，再次编译 ./build.sh

    #如何重启服务?
	scripts/control_web.sh status
    scripts/control_web.sh start
    scripts/control_web.sh stop
    scripts/control_web.sh restart

    # 如何查看数据?
    http://diabetes.chunyu.me/ivr/
    # 相关的nginx配置
	https://git.chunyu.me/health/nginx_conf/blob/feature/test/servers/locations/diabetes_urls.location
```

* 测试:
    * go get -u github.com/stretchr/testify/assert
    * 运行TestCase:
        *  go test git.chunyu.me/infra/fileupload/service -v -run "TestAudioOperation"
        * `如何跑TestCase?"


## 运维:
* 参考: http://beego.me/docs/install/bee.md
* bee pack 打包代码和编译结果
* 参考: http://beego.me/docs/deploy/
	* conf/app.conf
	* 这个部分如何定制呢?