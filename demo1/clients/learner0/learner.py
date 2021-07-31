import device
import multiprocessing
import os
import fldriver
import time


class Learner(device.Device):
    def __init__(self,name,edge_ip,edge_port,edge_name):
        super().__init__(name,edge_ip,edge_port,edge_name)

    def prepare_learning(self):
        if os.path.exists("stop"): 
            os.remove("stop")
        if os.path.exists("ready"):
            os.remove("ready")
        if os.path.exists("run"):
            os.remove("run")
        if os.path.exists("prepare"):
            os.remove("prepare")
        fldriver.text_create("prepare","prepare")
        
        
    def start_learning(self):
        fldriver.text_create("run","run")

    def stop_learning(self):
        fldriver.text_create("stop","stop")

    def check_task(self):
        info=self._get_twin_conf("taskinfo")
        if info=="ready":
            self.prepare_learning()
        elif info=="stop":
            self.stop_learning()
        elif info=="run":
            self.start_learning()

    def twin_change(self,obj):
        for x in obj:
            if obj[x]['desired'] != "":
                self._set_twin_conf(x,obj[x]['desired'])
        print("twin_change ok")
        self.check_task()
    
    def start_loop(self):
        while True:
            if os.path.exists("stop"): 
                os.remove("stop")
                if os.path.exists("ready"):
                    os.remove("ready")
                if os.path.exists("run"):
                    os.remove("run")
                self._set_twin_conf("status","Online")
                self._update_twin_direct()
            
            if os.path.exists("prepare"):
                os.remove("prepare")
                fl=fldriver.FL_Driver()
                mp1 = multiprocessing.Process(target=fl.ready_learning)
                mp1.start()

            if os.path.exists("ready") and self._get_twin_conf("status")!="ReadyForLearning":
                if self._get_twin_conf("status")!="Running":
                    self._set_twin_conf("status","ReadyForLearning")
                    self._update_twin_direct()

            if os.path.exists("run") and self._get_twin_conf("status")!="Running":
                self._set_twin_conf("status","Running")
                self._update_twin_direct()
                print("run!!!")
                fl=fldriver.FL_Driver()
                mp1 = multiprocessing.Process(target=fl.start_loop)
                mp1.start()
            
            time.sleep(1)
        
def main():
    l1=Learner("learner0","192.168.31.10",1883,"raspberrypi")
    res=l1.start_edge()
    l1.start_loop()

if __name__ == '__main__':
    main()
