FROM node:lts
WORKDIR /app
COPY twitcasting-manager/package*.json ./
RUN npm install --legacy-peer-deps
COPY twitcasting-manager/ ./
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
