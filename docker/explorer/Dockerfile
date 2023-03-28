# Dockerfile for Hyperledger Explorer Docker Image
FROM node:10.15-slim
MAINTAINER Datapace

# install maintain tools
RUN apt-get update
RUN apt-get install -y git bash python build-essential

# clone latest code from github
RUN git clone --single-branch -b v0.3.9.5 --depth 1 https://github.com/hyperledger/blockchain-explorer.git

WORKDIR /blockchain-explorer

RUN mkdir -p logs/app
RUN mkdir -p logs/db

# Modify config.json to update the value of pg host, username, password details.
# If you are building your code for production
RUN npm rebuild &&  npm install && cd client/ && npm install \
     && npm run build
EXPOSE 8080
CMD ["node", "main.js"]
