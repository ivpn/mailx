# build stage
FROM node:26.13.0-alpine3.24@sha256:c6ad9e9c1d80f6eb607c26b4b0a132ff7a2344fad1c78ea3baa2a7d4c8de6d94 AS build-stage
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
FROM nginxinc/nginx-unprivileged:1.31.6-alpine3.24@sha256:b9241c6e7b8e9a862f129d8d4199ab64b10390949a78bdd5603379b32c844083 AS production-build
USER root
COPY nginx.conf /etc/nginx/nginx.conf
RUN rm -rf /usr/share/nginx/html/*
COPY --from=build-stage /app/dist /usr/share/nginx/html/app
RUN chown -R nginx:nginx /usr/share/nginx/html
USER nginx
EXPOSE 3001
ENTRYPOINT ["nginx", "-g", "daemon off;"]
