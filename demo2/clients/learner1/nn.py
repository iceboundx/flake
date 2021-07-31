from collections import OrderedDict

import flwr as fl
import clientapp
import torch
import torch.nn as nn
import torch.nn.functional as F
import torchvision.transforms as transforms
from torch.utils.data import DataLoader
from torchvision.datasets import MNIST
import sys

DEVICE = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")


def main(url,key):
    """Create model, load data, define Flower client, start Flower client."""

    # Model (simple CNN adapted from 'PyTorch: A 60 Minute Blitz')
    class Net(nn.Module):
        def __init__(self):
            super(Net, self).__init__()
            self.layer1 = nn.Linear(28*28, 300)
            self.layer2 = nn.Linear(300, 100)
            self.layer3 = nn.Linear(100, 10)
    
        def forward(self, x):
            x = self.layer1(x)
            x = self.layer2(x)
            x = self.layer3(x)
            return x

    net = Net().to(DEVICE)

    # Load data (CIFAR-10)
    trainloader, testloader = load_data()

    # Flower client
    class CifarClient(fl.client.NumPyClient):
        def get_parameters(self):
            now= [val.cpu().numpy() for _, val in net.state_dict().items()]
            return now

        def set_parameters(self, parameters):
            params_dict = zip(net.state_dict().keys(), parameters)
            state_dict = OrderedDict({k: torch.Tensor(v) for k, v in params_dict})
            net.load_state_dict(state_dict, strict=True)

        def fit(self, parameters, config):
            self.set_parameters(parameters)
            train(net, trainloader, epochs=1)
            a=len(trainloader)
            return self.get_parameters(), len(trainloader)

        def evaluate(self, parameters, config):
            self.set_parameters(parameters)
            loss, accuracy = test(net, testloader)
            return len(testloader), float(loss), float(accuracy)

    # Start client
    print("start")
    clientapp.start_numpy_client(url, client=CifarClient(),keyPrefix=key)
    


def train(net, trainloader, epochs):
    """Train the network on the training set."""
    criterion = torch.nn.CrossEntropyLoss()
    optimizer = torch.optim.SGD(net.parameters(), lr=0.001, momentum=0.9)
    for _ in range(epochs):
        for images, labels in trainloader:
            img = images.view(images.size(0), -1)
            img, labels = img.to(DEVICE), labels.to(DEVICE)
            optimizer.zero_grad()
            loss = criterion(net(img), labels)
            loss.backward()
            optimizer.step()


def test(net, testloader):
    """Validate the network on the entire test set."""
    criterion = torch.nn.CrossEntropyLoss()
    correct, total, loss = 0, 0, 0.0
    with torch.no_grad():
        for data in testloader:
            images, labels = data[0].to(DEVICE), data[1].to(DEVICE)
            img = images.view(images.size(0), -1)
            outputs = net(img)
            loss += criterion(outputs, labels).item()
            _, predicted = torch.max(outputs.data, 1)
            total += labels.size(0)
            correct += (predicted == labels).sum().item()
    accuracy = correct / total
    return loss, accuracy


def load_data():
    """Load CIFAR-10 (training and test set)."""
    data_tf = transforms.Compose(
    [transforms.ToTensor(),
     transforms.Normalize([0.5], [0.5])])
    trainset = MNIST("data", train=True, download=True, transform=data_tf)
   # print(trainset.shape)
    testset = MNIST("data", train=False, download=True, transform=data_tf)
    trainloader = DataLoader(trainset, batch_size=32, shuffle=True)
    testloader = DataLoader(testset, batch_size=32)
    return trainloader, testloader


if __name__ == "__main__":
    url="flakeedge1.com:2333"
    key="edge1"
    if len(sys.argv)>1:
        url=sys.argv[1]
        key=sys.argv[2]
    main(url,key)
