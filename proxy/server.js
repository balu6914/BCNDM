const express = require('express');
const bodyParser = require('body-parser');
const app = express();
const bunyan = require('bunyan');
const proxy = new(require('./lib/proxy'));
// Instance API logger
const logger = bunyan.createLogger({name: 'api'});

app.use(bodyParser.json());
// Load API routes
require('./api')(app);

app.listen(8085);
logger.info('API is running on port 8085');

// Init Proxy server
proxy.setupHttpProxy();
