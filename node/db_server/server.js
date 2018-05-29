const http = require('http');
const express = require('express');
const bodyParser = require('body-parser');
const { Client } = require('pg');
const morgan = require('morgan');

const SERVER_PORT = 8080;
const LOGGER_FORMAT = ':method :url :status - :response-time ms';

const DB_HOST = 'localhost';
const DB_PORT = 5432;
const DB_NAME = 'instagram_development';
const DB_USER = 'postgres';
const DB_PASSWORD = 'postgres';

const postgres = new Client({
  host: DB_HOST,
  port: DB_PORT,
  database: DB_NAME,
  user: DB_USER,
  password: DB_PASSWORD,
});

(async () => {
  try {
    await postgres.connect();
  } catch (e) {
    console.error(e);
  }
})();

const app = express();
app.server = http.createServer(app);

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use(morgan(LOGGER_FORMAT));

app.get('/user/:id', async (req, res) => {
  const sql = 'SELECT * FROM users WHERE id = $1 LIMIT 1';
  const values = [req.params.id];

  try {
    const [user] = await postgres.query(sql, values);
    return res.json(user);
  } catch (e) {
    return res.status(500).json({ error: 'Internal Server Error' });
  }
});

app.server.listen(SERVER_PORT);
// eslint-disable-next-line no-console
console.log(`🚀  Server listening on port ${app.server.address().port}...`);
