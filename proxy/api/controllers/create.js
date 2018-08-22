const crypto = require('crypto');
const async = require('async');

const Store = new (require('../../lib/store'));
const Auth = new (require('../../lib/auth'));
const config = require('../../lib/config');
const bunyan = require('bunyan');
// Init module logger and set name
const logger = bunyan.createLogger({name: 'manager'});

/**
 * Create new proxy end point
 * @param  {[type]} req [description]
 * @param  {[type]} res [description]
 * @return {[Object]} req.body  [Requires ttl and url properties]
 */
exports.create = function(req, res) {
    // Check for required params
    if ((!req.body.url) || (!req.body.ttl)) {
        logger.error('Invalid request, required params missing');
        return res.status('401').send('Required parameters are not provided');
    }

    let url = req.body.url;
    let expiresIn = req.body.ttl;
    // Check if ttl is Int
    if (!Number.isInteger(expiresIn) || typeof url !== 'string') {
        logger.error('Invalid required params format.');
        return res.status('401').send('Required params is not in right format.');
    }
    /**
     * Runs the tasks array of functions in series,
     * each passing their results to the next in the array.
     */
    async.waterfall([
        generatePath,
        generateToken,
        storeRoute,
    ], (err, path, token) => {
        // Build a URL and return it as result
        let streamUrl = {
            url: `${config.proxyURL}${path}?token=${token}`,
        };
            if (err) {
                return res.status('401').send('Cant save data');
            }
            logger.info({stream_url: streamUrl}, 'New route registered');
            res.status('200').json(streamUrl);
        });
    /**
     * Generate unique stream route path
     * @param  {Function} done [Callback]
     */
    function generatePath(done) {
        crypto.randomBytes(16, (err, buffer) => {
            let path = buffer.toString('hex');
            if (err) return done(err);
            done(err, path);
        });
    }
    /**
     * Generate route access JWT token
     * @param  {String}   path [description]
     * @param  {Function} done [description]
     */
    function generateToken(path, done) {
        let data = {
            path: path,
            ttl: expiresIn,
        };
        Auth.createToken(data, (err, token) => {
            if (err) return done(err);
            done(err, path, token);
        });
    }
    /**
     * Store new proxy route in local storage
     * @param  {String}   path  [description]
     * @param  {String}   token [description]
     * @param  {Function} done  [description]
     */
    function storeRoute(path, token, done) {
        // Set the / prefix to path
        path = '/' + path;
        Store.set(path, url, (err, result) => {
            if (err) return done(err);
            done(err, path, token);
        });
    }
};
