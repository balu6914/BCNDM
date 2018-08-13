# Dockerfile for Hyperledger Explorer Docker Image
FROM node:9-slim
MAINTAINER Datapace

# install maintain tools
RUN apt-get update 
RUN apt-get install -y git bash python build-essential

# clone latest code from github
RUN git clone --single-branch -b release-3.4 --depth 1 https://github.com/hyperledger/blockchain-explorer.git

WORKDIR /blockchain-explorer

RUN mkdir -p logs/app 
RUN mkdir -p logs/db

# Modify config.json to update the value of pg host, username, password details.
# If you are building your code for production
# RUN npm install --only=production
RUN npm rebuild &&  npm install && cd client/ && npm install \
     && npm test -- -u --coverage && npm run build
#VOLUME /blockchain-explorer
EXPOSE 8080
CMD ["node", "main.js"]
