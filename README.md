# ARPScan--git
python简单实现ARP扫描 和 断网攻击工具, 采用单例设计模式, 方便调用.

# 安装说明
1. `pip install requirements.txt`
2. `python main.py`

# 使用
+ main.py 中填好自己的信息运行即可
+ 自己ip可以通过`ifconfig`获取,mac地址可以使用`getYourNetwork.py`获取
+ 启动后自动进入循环,新加入的mac地址都会被断网
+ 白名单中填写不想被断网的 ip

# 函数说明
- process_run(): 通过发送伪造的arp包来达到断网效果
- ARP_Scan.arpspoof(): 读取扫描到的客户端ip和mac，判断白名单和已经攻击的对象，给每一个客户端创建一个process进程来进行独立的攻击。
- ARP_Scan.send_Package(): 发送arp包来扫描存活设备'plf:AlienBlue_walk2'
