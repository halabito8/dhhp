FROM node:20-slim

WORKDIR /app

COPY entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

RUN npm install -g npm@latest create-vite vite

EXPOSE 5173

ENTRYPOINT ["entrypoint.sh"]