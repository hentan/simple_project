# ---------- build stage ----------
FROM node:20-alpine AS build
WORKDIR /app

COPY package.json ./
# For small demo projects, npm install is fine without lockfile
RUN npm install

COPY . .
RUN npm run build

# ---------- runtime stage ----------
FROM nginx:1.25-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
