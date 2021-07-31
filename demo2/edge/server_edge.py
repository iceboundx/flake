import flwr as fl
import sys
import app
import numpy as np
import os
from typing import Callable, Dict, List, Optional, Tuple
import subprocess
import time
k=0
p=None
def set_states(statename):
    file = open(statename, 'w+')
    file.write(statename)
    file.close()

def clear_states(statename):
    os.remove(statename)

def check_states(statename):
    time.sleep(0.3)
    return os.path.exists(statename) 
class SaveModelStrategy(fl.server.strategy.FedAvg):
    def __init__(self):
        super().__init__()
    
    def aggregate_fit(
        self,
        rnd: int,
        results: List[Tuple[fl.server.client_proxy.ClientProxy, fl.common.FitRes]],
        failures: List[BaseException],
    ) -> Optional[fl.common.Weights]:
        global k
        global p
        aggregated_weights = super().aggregate_fit(rnd, results, failures)
        if aggregated_weights is not None:
            # Save aggregated_weights
            print(f"Saving round {rnd} aggregated_weights...")
            np.savez(f"edge_round-{rnd}-weights.npz", *aggregated_weights)
            if rnd%k==0:
                np.savez(f"round-weights.npz", *aggregated_weights)
                with open(f'round-fitres.npz', 'w+') as f:
                    f.write(str(len(results))+'\n')
                    for x in results:
                        f.write(str(x[1].num_examples)+'\n')
                set_states("ready")
                while(1):
                    if check_states("aggok"):
                        break
                time.sleep(0.3)
                clear_states("aggok")
                test_data = np.load(f"finish_round-weights.npz")
                aggregated_weights = [val for _, val in test_data.items()]                
        return aggregated_weights



if __name__ == "__main__":
    k=3
    nums_rounds=30
    ipport="0.0.0.0:2333"
    cloudurl="flakeaggcore.com:6666"
    keyprefix="server"
    if len(sys.argv)>1:
        k=int(sys.argv[1])
        nums_rounds=int(sys.argv[2])
        ipport=sys.argv[3]
        cloudurl=sys.argv[4]
        keyprefix=sys.argv[5]
    strategy = SaveModelStrategy()
    p = subprocess.Popen("python3 client_edge.py "+cloudurl,shell=True,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
    app.start_server(ipport, strategy=strategy,config={"num_rounds": nums_rounds},key_prefix=keyprefix)
    