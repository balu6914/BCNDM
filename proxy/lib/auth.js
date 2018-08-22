/**
 * [Auth description]
 *
 */
const jwt = require('jsonwebtoken');
const config = require('./config');
const moment = require('moment');

/**
 * [Auth description]
 * @constructor
 */
class Auth {
    /**
     * [description]
     * @param  {Object} data [Data object is required, contains tll and stream route path]
     * @param  {Function} next [Callback]
     */
    createToken(data, next) {
        // Sign token, set exparation depending of
        // subscription durition (provided time to live)
        const token = jwt.sign({
            endDate: moment().add(data.ttl, 'hours'),
            path: data.path,
            issuer: config.issuerName,
        }, config.secret, {expiresIn: data.ttl + 'h'});
        next(null, token);
    }
    /**
     * [validateToken description]
     * @param  {[type]}   token [description]
     * @param  {Function} next  [Callback]
     */
    validateToken(token, next) {
        jwt.verify(token, config.secret, function(err, decoded) {
            if (err) {
                next(err, false);
            } else {
                next(null, true);
            }
        });
    }
}

module.exports = Auth;
