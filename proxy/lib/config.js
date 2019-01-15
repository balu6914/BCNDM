const appConfig = {
    secret: 'ynehj45le70cfbrcsszm', // JWT secret key
    issuerName: 'Datapace',
    proxyURL: process.env.PUBLIC_URL || 'http://localhost:8081', // Public access UR
    redisUrl: process.env.REDIS_URL || 'redis://localhost:6379',
};

module.exports = appConfig;
