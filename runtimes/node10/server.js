const morphModule = require('./index');
const express = require('express');
const app = express();

const endpoint = process.env.ENDPOINT || 'module-a';

// We need to create two endpoints, one which wraps around the render function
// a second which is a health-check endpoint. This endpoint is crucial.
app.get(`/${endpoint}`, (req, res) => res.send(morphModule.render()));
app.get('/_health', (req, res) => res.send('OK'));

app.listen(8080, () => console.log('Example app listening on port 8080!'));
