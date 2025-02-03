FROM node:lts
ARG NEXT_PUBLIC_MANAGE_BACKEND_ORIGIN
ENV NEXT_PUBLIC_MANAGE_BACKEND_ORIGIN=$NEXT_PUBLIC_MANAGE_BACKEND_ORIGIN
WORKDIR /app
COPY twitcasting-manager/package*.json ./
RUN npm install --legacy-peer-deps
COPY twitcasting-manager/ ./
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
