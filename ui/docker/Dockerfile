# Stage 0, based on Node.js, to build and compile Angular
FROM node:lts-alpine3.10 as node
WORKDIR /app
COPY package.json /app/
RUN npm install
COPY ./ /app/
RUN npm run build --prod

# Stage 1, based on Nginx, to have only the compiled app, ready for production with NginX
FROM nginx:alpine
COPY --from=node /app/dist/ /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
RUN rm -rf /usr/share/nginx/html/50x.html
