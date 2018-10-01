FROM node:8.9.4-alpine

MAINTAINER Nikola Marcetic <nikola@mainflux.com>

RUN mkdir /src

RUN npm config set unsafe-perm true
RUN npm install nodemon -g

WORKDIR /src
ADD package.json /src/package.json
RUN npm install

ADD nodemon.json /src/nodemon.json
# Copy source code
COPY . /src

EXPOSE 8085 8081

CMD npm run nodemon
