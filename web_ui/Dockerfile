FROM node:16.14.2-alpine

COPY ./public /home/app/public
COPY ./src /home/app/src
COPY package.json ./

WORKDIR /home/app

RUN npm install

CMD ["sh", "-c", "node ./src/main.js"]