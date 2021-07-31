from collections import OrderedDict
import mqttclient
import os
import nn
import subprocess
import time


def text_create(name, msg):
    file = open(name, 'w+')
    file.write(msg)
    file.close()


class FL_Driver:
    
    def __init__(self):
        self.p=None

    def ready_learning(self):
        text_create("ready","ready")
    
    def start_learning(self):
        print("learning start")
        self.p = subprocess.Popen("python nn.py",shell=True,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
        print(self.p)

    def stop_learning(self):
        if self.p is not None:
            self.p.kill()
    
    def check_learning(self):
        if os.path.exists("stop"):
            self.stop_learning()
            return False
        if self.p.poll() is not None:
            text_create("stop","stop")
            return False
        return True

    def start_loop(self):
        print("start loop")
        if os.path.exists("stop") or os.path.exists("ready")==False:
            print("Can not start loop")
            return
        self.start_learning()
        while self.check_learning():
            with open("fldriver.txt", "a", encoding="utf-8") as f:
                f.write(self.p.stdout.readline().decode("gbk").strip() + "\n") 
        print("process stop")
    

