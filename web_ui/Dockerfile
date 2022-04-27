FROM node:16.14.2-alpine

HEALTHCHECK NONE

COPY ./public /home/app/public
COPY ./src /home/app/src
COPY package.json ./

RUN useradd -ms /bin/bash nodeuser
USER nodeuser

WORKDIR /home/app

RUN npm update && npm install

CMD ["sh", "-c", "node ./src/main.js"]