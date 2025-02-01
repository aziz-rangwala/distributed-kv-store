FROM node:18-alpine

WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install

# Copy all source code
COPY . .

# Expose the port that the development server uses
EXPOSE 3000

# Start the React development server
CMD ["npm", "start"]
