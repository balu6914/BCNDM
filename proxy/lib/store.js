const redis = require('redis');
const config = require('./config');
const bunyan = require('bunyan');
// Init module logger and set logger name
const logger = bunyan.createLogger({name: 'redis'});

// Connect to redis
let client = redis;
client = redis.createClient(config.redisUrl);
// Listen to Redis connection error events
client.on('error', function(err) {
    logger.error('Can\'t connect to Redis!');
    process.exit();
});
/**
 * Store service is responsible for communication with redis, the main purpose
 * is to encapsulate the storage service and provide API interface.
 */
class Store {
    /**
     * [Set will add new route to routing table]
     * @param {String}   key   [Key to store]
     * @param {String}   value [Value to store]
     * @param {Function} next  [description]
     */
    set(key, value, next) {
        client.set(key, value, (err, data) => {
            if (err) next(err);
            // Emit the event, new value is saved in routing table
            next(err, data);
        });
    }
    /**
     * [Get will fetch route from routing table]
     * @param  {String}   key  [Key to get]
     * @param  {Function} next [Callback]
     */
    get(key, next) {
        client.get(key, function(err, reply) {
            if (err) return next(err);
            // NOTE: reply is null when the key is missing
            next(err, reply);
        });
    }
    /**
     * [Remove will remove route from routing table]
     * @param  {String}   key  [Key to remove]
     * @param  {Function} next [Callback]
     */
    remove(key, next) {
        client.del(key, function(err, reply) {
            if (err) return next(err);
            next(err, reply);
        });
    }
}

module.exports = Store;
