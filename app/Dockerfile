# build stage
FROM node:lts-alpine AS build-stage
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
EXPOSE 3001
ENTRYPOINT ["nginx", "-g", "daemon off;"]
