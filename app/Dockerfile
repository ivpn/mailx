# build stage
FROM node:21-alpine AS build-stage
WORKDIR /app
COPY package.json .
RUN npm install
COPY . .
RUN npm run build

# production stage
FROM nginx:alpine as production-build
COPY nginx.conf /etc/nginx/nginx.conf
RUN rm -rf /usr/share/nginx/html/*
COPY --from=build-stage /app/dist /usr/share/nginx/html/app
COPY ./access_rules.conf /etc/nginx/access_rules.conf
EXPOSE 3001
ENTRYPOINT ["nginx", "-g", "daemon off;"]
