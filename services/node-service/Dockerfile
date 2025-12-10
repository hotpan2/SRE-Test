FROM node:18-alpine
WORKDIR /app
COPY package.json .
RUN npm install --production
COPY . .
ENV PORT=8081
EXPOSE 8081
CMD ["node", "app.js"]
