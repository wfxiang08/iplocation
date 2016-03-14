# iplocation实现IP到地理位置的转换

## 1. 环境搭建:(采用gpm来实现类似Python Virtualenv的项目环境管理)

* `安装工具`
	* MAC OS: ```brew install gpm```
	* 其他: ```wget https://raw.githubusercontent.com/pote/gpm/v1.4.0/bin/gpm && chmod +x gpm && sudo mv gpm /usr/local/bin```

* 工作目录: 约定使用 `~/goprojects/`
	* 项目初建:
		* mkdir iplocation && cd iplocation
		* mkdir src && mkdir bin && mkdir pkg

* `下载代码`
	* cd ~/goprojects/
	* git clone git@git.chunyu.me:infra/iplocation.git
	* cd iplocation
	* source start_env.sh
	* cd src && gpm install (按照Godeps中的配置来布置环境)

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