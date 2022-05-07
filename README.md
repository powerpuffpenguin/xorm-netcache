# xormcache

[xorm.io](https://xorm.io/) is a database xorm library for golang. The [xorm.io](https://xorm.io/) library provides an interface to implement the cache and already provides a localized lru cache. But localized caches are not suitable for microservices, because there may be multiple different services in microservices that use the database so they need a common network cache. 

I wanted to use [redis](https://redis.io/) as a network cache for [xorm.io](https://xorm.io/), so I wrote this library


