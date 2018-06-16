FROM node:10

RUN mkdir -p /usr/app

WORKDIR /usr/app

ADD server.js server.js

ADD package.json package.json

RUN npm install

CMD ["node", "server.js"]
