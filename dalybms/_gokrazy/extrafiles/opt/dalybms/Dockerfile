FROM python:3.11-slim

RUN apt update && apt install -y python3-pip rfkill

RUN pip3 install pyserial bleak paho-mqtt

COPY python-daly-bms python-daly-bms
RUN cd python-daly-bms && python3 setup.py install

# setup startup script
COPY entrypoint.sh .
RUN chmod +x ./entrypoint.sh
CMD ./entrypoint.sh

