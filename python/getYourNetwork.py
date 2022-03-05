#!/usr/bin/python
# -*- coding: UTF-8 -*-
from scapy.all import *
print("YOUR Interfaces :")
print(show_interfaces())
print(getmacbyip("192.168.0.32"))

