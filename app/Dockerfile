# build stage
FROM node:22-alpine@sha256:c610fcdfb1d5b4740dd70c284ed3cb16bb857e0f7166196e36a5501df7a3aa32 AS build-stage
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile
COPY . .

# VITE_* build args are embedded into the compiled JS by `vite build` - app/.env
# is excluded from the build context (.dockerignore), so these must come in via
# --build-arg rather than a copied .env file.
ARG VITE_API_URL
ARG VITE_DOMAINS
ARG VITE_APP_NAME
ARG VITE_RESYNC_URL
ARG VITE_PRICING_URL
ENV VITE_API_URL=$VITE_API_URL \
    VITE_DOMAINS=$VITE_DOMAINS \
    VITE_APP_NAME=$VITE_APP_NAME \
    VITE_RESYNC_URL=$VITE_RESYNC_URL \
    VITE_PRICING_URL=$VITE_PRICING_URL

# fail loudly instead of silently baking an empty API URL into the bundle
RUN test -n "$VITE_API_URL" || (echo "ERROR: --build-arg VITE_API_URL is required" >&2; exit 1)

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
