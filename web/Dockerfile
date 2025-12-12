FROM node:20-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

EXPOSE 8000

# Run dev server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
