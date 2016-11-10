from nginx:alpine

RUN apk add -U nodejs

RUN npm install -g typescript
RUN npm install typings --global

COPY . /usr/share/nginx/html

WORKDIR /usr/share/nginx/html

RUN npm install
RUN npm run tsc
