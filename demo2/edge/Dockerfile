FROM waggle/plugin-torch:1.4.0
LABEL maintainers="icebound"
LABEL description="Flake Edge aggcore"
USER root
RUN mv /etc/apt/sources.list /etc/apt/sources.list.bak && \
	echo "deb http://mirrors.tuna.tsinghua.edu.cn/raspbian/raspbian/ buster main non-free contrib rpi" >/etc/apt/sources.list && \
    echo "deb-src http://mirrors.tuna.tsinghua.edu.cn/raspbian/raspbian/ buster main non-free contrib rpi" >>/etc/apt/sources.list
RUN gpg --keyserver keyserver.ubuntu.com --recv 9165938D90FDDD2E 
RUN gpg --export --armor 9165938D90FDDD2E | apt-key add -
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  build-essential \
  gcc
RUN apt-get install -y build-essential python-dev
COPY grpcio-1.31.0rc1-cp36-cp36m-linux_armv7l.whl ./
RUN pip3 install grpcio-1.31.0rc1-cp36-cp36m-linux_armv7l.whl
RUN pip3 install --upgrade Cython -i https://pypi.douban.com/simple
COPY requirements.txt ./
RUN pip3 install -r /requirements.txt -i https://pypi.douban.com/simple
# Copy from build directory
COPY *.py ./
COPY *.crt ./
COPY *.key ./
