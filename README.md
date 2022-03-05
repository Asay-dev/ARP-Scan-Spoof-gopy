# ARTScript-ARP

ARTScript
ARTCyber
ARTCoder

ARTScript提供的ARP工具包,包含ARP扫描攻击等各类功能,提供了多种语言的版本方便在不同平台使用.
**[ ! ]** 手机端TERMUX使用golang版本

+ golang 实现arpposion
+ python 简单实现ARP扫描 和 断网攻击工具, 采用单例设计模式, 方便调用.

# 安装说明

+ python
  1. `pip install requirements.txt`
  2. `python main.py`
+ golang
  1. `go mod init ""`
  2. `go mod tidy`

# 使用

+ main.py 中填好自己的信息运行即可
+ 自己ip可以通过`ifconfig`获取,mac地址可以使用`getYourNetwork.py`获取
+ 启动后自动进入循环,新加入的mac地址都会被断网
+ 白名单中填写不想被断网的 ip
