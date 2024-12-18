# Stage 1: Build the admin-panel frontend
FROM node:18-alpine AS admin-builder

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install

COPY . . 
RUN yarn build

# Stage 2: Serve the admin-panel using Nginx
FROM nginx:alpine

COPY --from=admin-builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]