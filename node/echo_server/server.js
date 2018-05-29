const http = require('http');
const express = require('express');
const bodyParser = require('body-parser');
const morgan = require('morgan');

const PORT = 8080;
const LOGGER_FORMAT = ':method :url :status - :response-time ms';

const app = express();
app.server = http.createServer(app);

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use(morgan(LOGGER_FORMAT));

app.get('/:action', (req, res) => res.json({ message: req.params.action }));

app.server.listen(PORT);
// eslint-disable-next-line no-console
console.log(`🚀  Server listening on port ${app.server.address().port}...`);
