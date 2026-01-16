FROM node:20
WORKDIR /app

COPY package.json ./
RUN npm install

COPY tsconfig.json ./
COPY src ./src
COPY retail.csv ./retail.csv

RUN npm run build

EXPOSE 8080
CMD ["npm", "start"]
