#!/bin/bash -x

DUSER=$(whoami)
DUID=${UID}

docker run --name stm32loader --rm -v $(pwd):/root/arduino-stm32-loader \
	-w /root/arduino-stm32-loader golang:latest \
	bash -c "useradd ${DUSER} -u ${DUID} && \
	apt-get update && apt-get install -y tar bzip2 git zip && \
	chmod +x *.sh && ./setup.sh && ./deploy-docker.sh && \
	chown ${DUSER}:${DUSER} -R distrib/ && \
	rm arduinoSTM32load && rm -rf bin/"
