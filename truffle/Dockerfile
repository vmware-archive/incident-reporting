# Copyright 2019 VMware, Inc.
# SPDX-License-Identifier: BSD-2
FROM node:chakracore-10.13.0

WORKDIR /usr/src/app

COPY truffle/package-lock.json .
COPY truffle/package.json .
RUN npm install && npm install -g truffle@4.1.14
COPY truffle/contracts contracts
COPY truffle/migrations migrations
COPY truffle/test test
COPY truffle/truffle.js .

CMD [ "bash" ]