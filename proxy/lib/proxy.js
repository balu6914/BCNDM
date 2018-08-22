const http = require('http');
const httpProxy = require('http-proxy');
const url = require('url');
const Auth = new (require('./auth'));
const Store = new (require('./store'));
const bunyan = require('bunyan');
const logger = bunyan.createLogger({name: 'proxy'});

/**
 * Proxy module is responsible for proxying requests to origin target,
 * which is saved in local storage.
 */
class Proxy {
    /**
     * [constructor description]
     * @param {[type]} opts [description]
     */
    constructor(opts) {
        if (!(this instanceof Proxy)) {
            return new Proxy(opts);
        }
        this.opts = opts || {};
        this.logger = logger || bunyan.createLogger({name: 'proxy'});
    }
    /**
     * [setupHttpProxy description]
     * @param  {[type]} logs        [description]
     * @return {Object}             [New Proxy server instance]
     */
    setupHttpProxy() {
        // Set up proxy rules instance
        http.createServer((req, res) => {
            // Athorize request  and validate token
            this.authorize(req, res, (err) => {
                if (err) return this.unauthorized(req, res);
                // Check target route
                this.findRoute(req, res, (err, url) => {
                    if (err) return this.notFound(url, req, res);
                    const proxy = httpProxy.createProxyServer({});
                    // Proxy request
                    proxy.web(req, res, {
                        target: url,
                        secure: true, // Depends on your needs, could be false.
                        // ws: true // WS support
                        ignorePath: true,
                        changeOrigin: true
                    });
                    // Set req headers
                    proxy.on('proxyReq', (p, req) => {
                        if (req.host != null) {
                            p.setHeader('host', req.host);
                        }
                    });
                    // Listen for the `error` event on `proxy,
                    // use case when the Stream provider/API is down`.
                    proxy.on('error', (err, req, res) => {
                        this.error(err, req, res);
                    });
                });
            });
        })
        .listen(8081);
        this.logger.info('Proxy is running on port 8081');
    }
    /**
     * Remove route when subscription expires
     * @param  {[type]} path [description]
     */
    removeRoute(path, next) {
        Store.remove(path, (err, result) => {
            if (err) return next(err);
            this.logger.info('Route removed', {proxy_path: path})
            next(err, path);
        });
    }
    /**
     * Find route target in local storage
     * @param  {[type]}   req  [description]
     * @param  {[type]}   res  [description]
     * @param  {Function} next [description]
     */
    findRoute(req, res, next) {
        // Get route path from request URL
        let path = url.parse(req.url, true).pathname || '';
            Store.get(path, (err, result) => {
                let notFound = 'Route path does not exist.';
                if (err) return next(err, path);
                // NOTE: If not result store service returns null (redis way)
                if (!result) return next(notFound, path);
                next(err, result);
            });
    }
    /**
     * Route does not exist
     * @param  {Object} path [description]
     * @param  {Object} req [description]
     * @param  {Object} res [description]
     */
    notFound(path, req, res) {
        res.statusCode = 404;
        res.write('Not Found');
        res.end();
        this.logger.error({proxy_path: path}, 'Route path not found');
    }
    /**
     * Validate JWT token and subscription validity
     * @param  {Object}   req  [description]
     * @param  {Object}   res  [description]
     * @param  {Function} next [description]
     */
    authorize(req, res, next) {
        let query = url.parse(req.url, true).query;
        if (!query.token) {
            let err = 'No token provided';
            next(err);
        }
        Auth.validateToken(query.token, (err) => {
            // Check if token is exipired
            if(err && err.name === 'TokenExpiredError') {
                // Remove route, subscription is not a valid.
                let path = url.parse(req.url, true).pathname || '';
                this.removeRoute(path, (error) => {
                    if(error) return next(error);
                    next(err)
                })
            }
            next(err);
        });
    }
    /**
     * unauthorized access handler
     * @param  {Object} req [description]
     * @param  {Object} res [description]
     */
    unauthorized(req, res) {
        res.statusCode = 403;
        res.end('Access rejected! Your authorization is not valid or your subscription is exipired.');
        this.logger.error('Unauthorized');
    }
    /**
     * Proxy response error handler
     * @param  {Object} err [description]
     * @param  {Object} req [description]
     * @param  {Object} res [description]
     */
    error(err, req, res) {
        let path = url.parse(req.url, true).pathname || '';
        res.statusCode = 503;
        res.end('Something went wrong. Stream provider is probably down.');
        this.logger.error({proxy_path: path, proxy_error: err.reason || {}}, 'Cant access the stream!');
    }
}

module.exports = Proxy;
