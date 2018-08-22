const createController = require('./controllers/create');
const bodyParser = require('body-parser');

const bodyParserMiddlewares = [
    bodyParser.urlencoded({extended: false}),
    bodyParser.json({limit: '16mb'}),
];

module.exports = (app) => {
    app.post('/api/register', bodyParserMiddlewares, createController.create);
};
