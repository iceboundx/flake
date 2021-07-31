import mqttclient
import time
import math
import json

class twinctl:
    def connect_ok(self,client, userdata, flags, rc):
        print("Connected with result code: " + str(rc))
        self._connect_flag=1

    def subscribe_ok(self,client, userdata, mid, granted_qos):
        print("On Subscribed: qos = %d" % granted_qos)
        self._subscribe_cnt=self._subscribe_cnt+1

    def messagecontrol(self,topic,obj):
        print("message get from "+topic+": "+str(obj))
        if "code" in obj:
            if obj['code']==404:
                print("message not found error")
                return
        if "error" in obj:
            print("message control error")
            return
        if topic==self._get_topic_res:
            self._get_topic_res_callback(obj)
        elif topic==self._update_topic_res:
            self._update_res_callback(obj)
        elif topic==self._twin_change_topic:
            self._twin_change_callback(obj)
        else:
            print("no match topic:"+topic)

    def _update_res_callback(self,obj):
        now_time_stamp=int(time.time())
        if abs(obj['timestamp']/1000-now_time_stamp)<100:
            self._update_flag=0
    
    def _get_topic_res_callback(self,obj):
        now_time_stamp=int(time.time())
        if abs(obj['timestamp']/1000-now_time_stamp)<100:
            self._get_flag=0
            self.get_twin=dict()
            if 'twin' not in obj:
                return
            for x in obj['twin']:
                desired=""
                reported=""
                now_obj=obj['twin'][x]
                if "actual" in now_obj:
                    if "value" in now_obj["actual"]:
                        reported=now_obj["actual"]["value"]

                if "expected" in now_obj:
                    if "value" in now_obj["expected"]:
                        desired=now_obj["expected"]["value"]

                self.get_twin[x]={"desired":desired,"reported":reported}

    def _twin_change_callback(self,obj):
        now_time_stamp=int(time.time())
        if "event_id" not in obj:
            print("no event_id, not cloud, return")
            return
        if obj['event_id'] =='':
            print("no event_id, not cloud, return")
            return
        if abs(obj['timestamp']/1000-now_time_stamp)<100:
            self._change_flag=1
            self.get_twin=dict()
            if 'twin' not in obj:
                return
            for x in obj['twin']:
                desired=""
                reported=""
                if "current" not in obj['twin'][x]:
                    continue 
                now_obj=obj['twin'][x]["current"]
                if now_obj is None:
                    continue
                if "actual" in now_obj:
                    if "value" in now_obj["actual"]:
                        reported=now_obj["actual"]["value"]

                if "expected" in now_obj:
                    if "value" in now_obj["expected"]:
                        desired=now_obj["expected"]["value"]

                self.get_twin[x]={"desired":desired,"reported":reported}
            self.twin_change_callback(self.get_twin)

    def twin_change_callback(self,obj):#need override
        self._change_flag=0

    def get_twins(self):
        self._get_flag=1
        self._mqttcl.send(self._get_topic,self._get_str)
        cnt=0
        while self._get_flag and cnt<self._time_long:
            time.sleep(1)
            cnt=cnt+1
        if cnt>=self._time_long:
            return False
        return True
    
    def update_twins(self,update_dict):
        if isinstance(update_dict,dict)==False:
            return False
        final_dict=dict()
        print("update_twins:"+str(update_dict))
        for x in update_dict:
            now_dict={
                "actual":{
                    "value":update_dict[x]
                },
                "metadata":{"type":"string"},
            }
            final_dict[x]=now_dict
        zip_dict={
            "event_id":"",
            "timestamp":int(time.time()*1000),
            "twin":final_dict,
        }
        change_json=json.dumps(zip_dict)
        self._mqttcl.send(self._update_topic,change_json)
        cnt=0
        self._update_flag=1
        while self._update_flag and cnt<self._time_long:
            time.sleep(1)
            cnt=cnt+1
        if cnt>=self._time_long:
            return False
        return True

    def update_twins_direct(self,update_dict):
        if isinstance(update_dict,dict)==False:
            return False
        final_dict=dict()
        print("update_twins:"+str(update_dict))
        for x in update_dict:
            now_dict={
                "actual":{
                    "value":update_dict[x]
                },
                "metadata":{"type":"string"},
            }
            final_dict[x]=now_dict
        zip_dict={
            "event_id":"",
            "timestamp":int(time.time()*1000),
            "twin":final_dict,
        }
        change_json=json.dumps(zip_dict)
        self._mqttcl.send(self._update_topic,change_json)
        return True

    def __init__(self,name,ip,port):
        self._mqttcl=mqttclient.MqttClient(ip,port,self.messagecontrol)
        self._mqttcl.client.on_connect=self.connect_ok

        self.name=name
        self.get_twin=None

        self._get_str = '{"event_id":"","timestamp":0}'
        self._get_topic="$hw/events/device/"+name+"/twin/get"
        self._get_topic_res=self._get_topic+"/result"
        self._update_topic="$hw/events/device/"+name+"/twin/update"
        self._update_topic_res=self._update_topic+"/result"
        self._twin_change_topic=self._update_topic+"/document"

        self._time_long=10
        
        self._get_flag=0
        self._update_flag=0
        self._change_flag=0
        self._connect_flag=0
        self._subscribe_cnt=0

        self._all_subscribe=3
        self._mqttcl.subscribe(self._get_topic_res)
        self._mqttcl.subscribe(self._update_topic_res)
        self._mqttcl.subscribe(self._twin_change_topic)


    def start_client(self):
        self._mqttcl.start_loop()
        print("Init client, please wait.")
        nowcnt=0
        while self._connect_flag==False and self._subscribe_cnt<self._all_subscribe:
            time.sleep(1)
            nowcnt=nowcnt+1
            if nowcnt>15:
                return False
        print("Client init ok")
        return True

    def check_device(self):
        return self.get_twins()
        
if __name__ == '__main__':
    tw=twinctl("newlearner123","192.168.31.10",1883)
    err=tw.start_client()
    if err==False:
        print("Failed")
    res=tw.get_twins()
    print(tw.get_twin)
    if res==True:
        print("xxxx")
        print(tw.get_twin['status'])
    print("")
        
