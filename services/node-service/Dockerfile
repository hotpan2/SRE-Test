FROM node:18-alpine AS build
WORKDIR /app
COPY package.json .
RUN npm install

FROM node:18-alpine
WORKDIR /app
COPY --from=build /app/node_modules ./node_modules
COPY . .
ENV PORT=8081
EXPOSE 8081
CMD ["node", "app.js"]
