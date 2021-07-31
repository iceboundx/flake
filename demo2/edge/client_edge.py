from collections import OrderedDict

import flwr as fl
import numpy as np
import clientapp
import os
import time
import sys
from flwr.common.logger import log

param=None
num_examples=0
rounds=0

def check_states(statename):
    time.sleep(0.3)
    return os.path.exists(statename)

def clear_states(statename):
    os.remove(statename)

def set_states(statename):
    file = open(statename, 'w+')
    file.write(statename)
    file.close()


def get_from_file(nowround):
    global param
    global num_examples
    global rounds
    while(True):
        if check_states("ready"):
            break
    time.sleep(0.2)
    clear_states("ready")
    print("getget!!")
    weightname="round-weights.npz"
    fitresname="round-fitres.npz"
    test_data = np.load(weightname)
    param= [val for _, val in test_data.items()]
    with open(f'round-fitres.npz', 'r+') as f:
        num=int(f.readline())
        num_examples=0
        for i in range(num):
            num_examples+=int(f.readline())


def main(cloudurl):
    class CifarClient(fl.client.NumPyClient):

        def get_parameters(self):
            global param
            global num_examples
            global rounds
            print("get_param")
            if param==None:
                get_from_file(1)
                set_states("ready") #remain state
            return param

        def set_parameters(self, parameters):
            global param
            global num_examples
            global rounds
            print("set_param")
            param=parameters

        def fit(self, parameters, config):
            global param
            global num_examples
            global rounds
            print("fit")
            self.set_parameters(parameters)
            rounds=rounds+1
            get_from_file(rounds)
            return self.get_parameters(), num_examples

        def evaluate(self, parameters, config):
            print("evaluate param")
            global param
            global num_examples
            global rounds
            np.savez(f"finish_round-weights.npz",*parameters)
            set_states("aggok")
            return 100, 99.9, 0.9
        
    # Start client
    print("start")
    clientapp.start_numpy_client(cloudurl, client=CifarClient())
    
    

if __name__ == "__main__":
    cloudurl="flakecloud.com:6666"
    if len(sys.argv)>1:
        cloudurl=sys.argv[1]
    main(cloudurl)



