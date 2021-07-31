try:
    import paho.mqtt.client as mqtt
except ImportError:
    print("MQTT client not find. Please install as follow:")
    print("git clone http://git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.python.git")
    print("cd org.eclipse.paho.mqtt.python")
    print("sudo python setup.py install")

import json
import asyncio


class MqttClient:
    def on_connect(self,client, userdata, flags, rc):
        print("Connected with result code: " + str(rc))

    def on_message(self,client, userdata, msg):
        #print(msg.topic + " " + str(msg.payload))
        try:
            x=json.loads(str(msg.payload)[2:][:-1])
        except:
            print("msg load error")
            self.message_control(msg.topic,{"error":msg.payload})
        else:
            self.message_control(msg.topic,x)

    def on_subscribe(self,client, userdata, mid, granted_qos):
        print("On Subscribed: qos = %d" % granted_qos)

    def on_disconnect(self,client, userdata, rc):
        if rc != 0:
            print("Unexpected disconnection %s" % rc)
    
    def message_control(self):
        return self._on_message

    def message_control(self, func):
        self._message_control=func

            
    client = mqtt.Client()
    
    def __init__(self,ip,port,messagecontrol):
        self.client.on_connect = self.on_connect
        self.client.on_message = self.on_message
        self.client.on_subscribe = self.on_subscribe
        self.client.on_disconnect = self.on_disconnect
        self.client.connect(ip, port, 600)
        self._message_control = None
        self.message_control=messagecontrol
    
    def send(self,topic,obj):
        print("publish: "+topic+" "+obj)
        messageInfo=self.client.publish(topic, payload=obj, qos=0)
        
    def subscribe(self, topic):
        print("subscribe: ",topic)
        self.client.subscribe(topic,qos=1)
    
    def start_loop(self):
        self.client.loop_start()

if __name__ == '__main__':
    cl=MqttClient("192.168.31.10",1883,messagectl)
    #cl.subscribe("$hw/events/device/learner456/twin/update/document")
    #cl.subscribe("$hw/events/device/learner456/twin/get/result")
    cl.subscribe("$hw/events/device/learner456/twin/update/result")
    
    cl.start_loop()
    #cl.send("$hw/events/device/counter/twin/update", '{"event_id":"","timestamp":0,"twin":{"status":{"actual":{"value":"8"},"metadata":{"type":"Updated"}}}}')
    while True:
        x=1

