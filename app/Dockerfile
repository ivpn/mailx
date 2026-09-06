# build stage
FROM node:22-alpine@sha256:c610fcdfb1d5b4740dd70c284ed3cb16bb857e0f7166196e36a5501df7a3aa32 AS build-stage
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile
COPY . .
RUN yarn build

# production stage - nginx-unprivileged already runs as a non-root user by design
FROM nginxinc/nginx-unprivileged:alpine@sha256:aa8c9087d36d93e9d650c5365f883b421e8214aedbad24ade52b844c583358f1 AS production-build
USER root
COPY nginx.conf /etc/nginx/nginx.conf
RUN rm -rf /usr/share/nginx/html/*
COPY --from=build-stage /app/dist /usr/share/nginx/html/app
RUN chown -R nginx:nginx /usr/share/nginx/html
USER nginx
EXPOSE 3001
ENTRYPOINT ["nginx", "-g", "daemon off;"]
