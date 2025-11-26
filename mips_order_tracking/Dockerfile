#### STAGE 1: Build the React App ####

FROM node:20-alpine AS build
WORKDIR /app

ARG PUBLIC_API_URL
ARG PUBLIC_TOMTOM_API_KEY
ENV PUBLIC_API_URL=$PUBLIC_API_URL
ENV PUBLIC_TOMTOM_API_KEY=$PUBLIC_TOMTOM_API_KEY

COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

#### STAGE 2: Serve the App with Nginx ####

# Use a tiny, secure Nginx image
FROM nginx:stable-alpine

# Copy the static files from the build stage
COPY --from=build /app/dist /usr/share/nginx/html

# Copy the custom Nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80 (Nginx's default)
EXPOSE 80

# Run Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]