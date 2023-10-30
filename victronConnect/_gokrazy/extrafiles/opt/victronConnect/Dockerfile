FROM python:3.7-slim

RUN apt update && apt install -y python3-pip libcairo2-dev python3-venv python3-wheel python3-dev libgirepository1.0-dev build-essential python3-dbus libdbus-glib-1-dev

COPY victron /victron
WORKDIR /victron
RUN pip3 install -r requirements.txt

# setup startup script
COPY entrypoint.sh .
RUN chmod +x ./entrypoint.sh
CMD ./entrypoint.sh

