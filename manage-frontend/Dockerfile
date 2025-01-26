# Use the official Node.js image as the base image
FROM node:lts

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY twitcasting-manager/package*.json ./

# Install dependencies
RUN npm install --legacy-peer-deps

# Copy the rest of the application code
COPY twitcasting-manager/ ./

# Build the application
RUN npm run build

# Expose the port the app runs on
EXPOSE 3000

# Start the application
CMD ["npm", "start"]
