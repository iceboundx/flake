client: general client codes
server: server codes (golang)
demo1: demo1 codes
demo2: demo2 codes
demo-model2: demo-model2 codes
resources: contains the network configuration and a whl file

Operating environment requirements:
KubeEdge>=1.5
Python>=3.6
Golang>=1.14

When making Docker image, you need to copy the grpcio-1.31.0rc1-cp36-cp36m-linux_armv7l in Resources to the same directory as Dockerfile.

I have modified some of the content of flower library about network connection, please pay attention to distinguish.

To match your own domain name, you may need to generate your own key file.

