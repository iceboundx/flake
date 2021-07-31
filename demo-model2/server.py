import flwr as fl
import sys
import app
import numpy as np
import os
from typing import Callable, Dict, List, Optional, Tuple



class SaveModelStrategy(fl.server.strategy.FedAvg):
    def __init__(self,taskId):
        super().__init__()
        self._taskId=taskId
        self._taskPath='task_'+taskId
    
    def aggregate_fit(
        self,
        rnd: int,
        results: List[Tuple[fl.server.client_proxy.ClientProxy, fl.common.FitRes]],
        failures: List[BaseException],
    ) -> Optional[fl.common.Weights]:
        aggregated_weights = super().aggregate_fit(rnd, results, failures)
        if aggregated_weights is not None:
            # Save aggregated_weights
            print(f"Saving round {rnd} aggregated_weights...")
            if not os.path.exists(self._taskPath):
                os.mkdir(self._taskPath)
            np.savez(self._taskPath+f"/round-{rnd}-weights.npz", *aggregated_weights)
        return aggregated_weights



if __name__ == "__main__":
    task_id="1122"
    if len(sys.argv)>1:
        task_id=sys.argv[1]
    strategy = SaveModelStrategy(task_id)
    #app.start_server("0.0.0.0:2333",config={"num_rounds": 3})
    app.start_server("0.0.0.0:2333", strategy=strategy,config={"num_rounds": 3})
    