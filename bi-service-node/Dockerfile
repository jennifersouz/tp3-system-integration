FROM node:20
WORKDIR /app
COPY package.json ./
RUN npm install
COPY src ./src
EXPOSE 8082
CMD ["npm", "start"]
