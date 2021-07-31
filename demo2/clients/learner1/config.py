import json
import os
import time

class Config:

    def _read_conf(self):
        if os.path.exists(self._path)==False:
            f=open(self._path,"w")
            f.write(json.dumps({
                "create_time":int(time.time()*1000),
            }))
            f.close()
        
        with open(self._path,encoding='utf-8') as f1:
            try:
                content=f1.read()
                self._conf=json.loads(content)
            except Exception as e:
                print("config read error: "+str(e))
                return False
            else:
                return True
    
    def _update_conf(self):
        with open(self._path,encoding='utf-8',mode='w') as f1:
            try:
                f1.write(json.dumps(self._conf))
            except Exception as e:
                print("config read error: "+str(e))
                return False
            else:
                return True

    def get_conf(self,conf_key):
        if self.check_ready()==False:
            print("config is not ready")
            return False
        if conf_key not in self._conf:
            return ""
        return self._conf[conf_key]
    
    def change_conf(self,conf_key,conf_val):
        if self.check_ready()==False:
            print("config is not ready")
            return False
        self._conf[conf_key]=conf_val
        return self._update_conf()

    def check_ready(self):
        return self._is_ready

    def __init__(self,conf_path):
        self._conf=None
        self._path=conf_path
        self._is_ready=self._read_conf()

if __name__ == '__main__':
    conf=Config("test.conf")
    #print(conf.get_conf("create_time"))
    #conf.change_conf("test_key","value")
    print(conf.get_conf("test_key"))
