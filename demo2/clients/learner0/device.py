import mqttclient
import twin
import config
import requests
import time
import asyncio
import _thread

DeviceStatusKey="status"
DeviceStatusNotReady         = "NotReady"
DeviceStatusOnline           = "Online"
DeviceStatusOffline          = "Offline"
DeviceStatusReadyForLearning = "ReadyForLearning"
DeviceStatusLearning         = "Learning"

FirstConnectTimeKey="FirstConnectTime"
ConnectUrl="https://flakecloud.com:443/connect"

class Device:

    def first_connect(self,url):
        print("Start First Connection")
        key=self._key.get_conf("key")
        postdatas = {'id': self._twin.name,'node':self._edge_name, 'key':key}
        res = requests.post(url, data=postdatas,cert=('server.crt','server.key'),verify=False)
        if res.text=="Connect Ok":
            return True
        return False
    
    def re_connect(self,url):
        print("Start Reconnection")
        key=self._key.get_conf("key")
        postdatas = {'id': self._twin.name,'node':self._edge_name, 'key':key}
        res = requests.post(url, data=postdatas,cert=('server.crt','server.key'), verify=False)
        if res.text=="Connect Ok":
            res=self._conf.change_conf(FirstConnectTimeKey,int(time.time()*1000))
            if res==False:
                print("Fail to save reconnection conf")
                return False
            time.sleep(10)
            print("Wait 10 sec for cloud-edge synchronization")
        return True

    def _set_twin_conf(self,key,value):
        self._twin_content[key]=value
        return self._conf.change_conf("twin_"+key,value)     

    def _get_twin_conf(self,key):
        x=self._conf.get_conf("twin_"+key)
        self._twin_content[key]=x
        return x
    
    def _get_all_twins_from_conf(self):
        for x in self._twin_list:
            self._get_twin_conf(x)

    def _fetch_twin(self):
        flag=self._twin.get_twins()
        if flag==False:
            print("Get Twin from edge error")
            return False
        twins=self._twin.get_twin
        for x in twins:
            if twins[x]['desired'] != "":
                self._set_twin_conf(x,twins[x]['desired'])
        return True
    
    def _update_twin(self):
        update_dict=dict()
        for x in self._twin_list:
            v=self._get_twin_conf(x)
            update_dict[x]=v
        return self._twin.update_twins(update_dict)
    
    def _update_twin_direct(self):
        update_dict=dict()
        for x in self._twin_list:
            v=self._get_twin_conf(x)
            update_dict[x]=v
        return self._twin.update_twins_direct(update_dict)
    
    def sync_twin(self):
        if self._fetch_twin()==True:
            return self._update_twin()
        return False

    def get_twin_content(self,key):
        if key not in self._twin_content:
            return ""
        return self._twin_content[key]
    
    def check_edge_ready(self):
        return self._edge_ok

    def set_status(self,value):
        return self._set_twin_conf(DeviceStatusKey,value)
       

    def __init__(self,name,edge_ip,edge_port,edge_name):
        
        self._name=name
        self._edge_name=edge_name

        self._conf=config.Config(name+"_conf.json")
        self._key=config.Config(name+"_key.json")
        self._twin_content=dict()

        self._twin_list=['status','deviceinfo','taskinfo']
        self._edge_ok=False

        if self._conf.check_ready() == False:
            print("conf not ready.")
            return

        self._get_all_twins_from_conf()
        self._twin=twin.twinctl(name,edge_ip,edge_port)
        self._twin.twin_change_callback=self.twin_change
        if self._conf.get_conf(FirstConnectTimeKey)=="":
            self.set_status(DeviceStatusNotReady)
        else:
            self.set_status(DeviceStatusOffline)
        

    def start_edge(self):
        if self._conf.get_conf(FirstConnectTimeKey)=="":
            res=self.first_connect(ConnectUrl)
            if res==False:
                print("Fail to first connect")
                return False
            res=self._conf.change_conf(FirstConnectTimeKey,int(time.time()*1000))
            if res==False:
                print("Fail to save first connection conf")
                return False
            time.sleep(10)
            print("Wait 10 sec for cloud-edge synchronization")

        self._edge_ok=self._twin.start_client()
        if self._edge_ok==False:
            print("Edge is not ready")
            self.set_status(DeviceStatusOffline)
            return False
        else:
            check_flag=self._twin.check_device()
            if check_flag==False:
                self.re_connect(ConnectUrl)
                check_flag=self._twin.check_device()
            
            if check_flag==False:
                print("Sync failed")
                self.set_status(DeviceStatusOffline)
                return False
            else:
                self.set_status(DeviceStatusOnline)
                print("Sycn Successed")
                self.sync_twin()
                return True

    def twin_change(self,obj):
        for x in obj:
            if obj[x]['desired'] != "":
                self._set_twin_conf(x,obj[x]['desired'])
        print("twin_change ok")
        self._update_twin_direct()

def main():
    l1=Device("learner1","192.168.31.10",1883,"raspberrypi")
    res=l1.start_edge()
    print(res)
    while True:
        x=input()
        if x=="ttt":
            l1.set_status("ReadyForLearning")
        elif x=="rrr":
            l1.sync_twin()

def test_https():
    print("Start First Connection")
    postdatas = {'id': "888",'node':"999"}
    res = requests.post(ConnectUrl, data=postdatas,verify='access.crt')
    print(res.text)

if __name__ == '__main__':
    #main()
    test_https()
        

