# ---------- Build Stage ----------
FROM node:22-alpine AS builder

WORKDIR /app

# Configure npm network resilience parameters to handle intermittent network timeouts
RUN npm config set fetch-retries 5 \
 && npm config set fetch-retry-mintimeout 20000 \
 && npm config set fetch-retry-maxtimeout 120000 \
 && npm config set fetch-timeout 300000

COPY package*.json ./

# Utilize BuildKit cache mount for npm cache and run npm install with retry protection
RUN --mount=type=cache,target=/root/.npm \
    npm install --prefer-offline --no-audit

COPY . .

RUN npm run build

# ---------- Runtime Stage ----------
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]