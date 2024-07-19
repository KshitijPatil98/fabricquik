FROM node:14-alpine
WORKDIR /${CC_NAME}
COPY . .
RUN npm install
EXPOSE ${CC_PORT}
CMD ["npm","start"]
