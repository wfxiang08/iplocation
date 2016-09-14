# 糖尿病400电话服务

## 开发环境
* 工作目录: 约定使用 `~/goprojects/`
	* mkdir -p iplocation/src/git.chunyu.me/infra/
	* cd `iplocation/src/git.chunyu.me/infra/`
	* git clone git@git.chunyu.me:infra/iplocation.git
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
			    * 输入`source start_env.sh`脚本中输出的地址，例如:
				    * /Users/feiwang/goprojects/iplocation


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


## 运维部署
* 编译:
	* cd ~/goprojects/iplocation/src/git.chunyu.me/infra/iplocation
	* 运行编译脚本:
		* `bash scripts/build_ip.sh`
