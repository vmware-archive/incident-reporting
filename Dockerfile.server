FROM node:chakracore-10.13.0

WORKDIR /usr/src/app

COPY package.json package.json
COPY client/package.json ./client/package.json

RUN npm install && \
  cd client && \
  npm install

COPY . .

EXPOSE 3000
EXPOSE 8545
CMD [ "./entrypoint.sh" ]
