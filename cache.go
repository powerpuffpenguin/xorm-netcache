package xormcache

import "strings"

type Cache struct {
	opts  *options
	store Store
}

func New(store Store, opt ...Option) (c *Cache, e error) {
	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	c = &Cache{
		store: store,
		opts:  &opts,
	}
	return
}
func (c *Cache) sqlPrefix(tableName string) string {
	return strings.Join(
		[]string{
			c.opts.prefix,
			`sql`,
			tableName + c.opts.sep,
		},
		c.opts.sep,
	)
}
func (c *Cache) sqlKey(tableName, sql string) string {
	return strings.Join(
		[]string{
			c.opts.prefix,
			`sql`,
			tableName,
			sql,
		},
		c.opts.sep,
	)
}
func (c *Cache) beanPrefix(tableName string) string {
	return strings.Join(
		[]string{
			c.opts.prefix,
			`bean`,
			tableName + c.opts.sep,
		},
		c.opts.sep,
	)
}
func (c *Cache) beanKey(tableName, id string) string {
	return strings.Join(
		[]string{
			c.opts.prefix,
			`bean`,
			tableName,
			id,
		},
		c.opts.sep,
	)
}
func (c *Cache) GetIds(tableName, sql string) interface{} {
	key := c.sqlKey(tableName, sql)
	b, e := c.store.Get(key)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("GetIds(%s,%s) error: %s\n", tableName, sql, e)
		}
		return nil
	} else if b == nil {
		return nil
	}
	var result interface{}
	e = c.opts.coder.Decode(b, &result)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("GetIds(%s,%s) error: %s\n", tableName, sql, e)
		}
		return nil
	}
	return result
}
func (c *Cache) PutIds(tableName, sql string, ids interface{}) {
	value, e := c.opts.coder.Encode(ids)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("PutIds(%s,%s,%v) error: %s\n", tableName, sql, ids, e)
		}
		return
	}
	key := c.sqlKey(tableName, sql)
	e = c.store.Put(key, value)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("PutIds(%s,%s,%v) error: %s\n", tableName, sql, ids, e)
		}
		return
	}
}
func (c *Cache) DelIds(tableName, sql string) {
	key := c.sqlKey(tableName, sql)
	e := c.store.Del(key)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("DelIds(%s,%s) error: %s\n", tableName, sql, e)
		}
		return
	}
}

func (c *Cache) GetBean(tableName string, id string) interface{} {
	key := c.beanKey(tableName, id)
	b, e := c.store.Get(key)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("GetBean(%s,%s) error: %s\n", tableName, id, e)
		}
		return nil
	} else if b == nil {
		return nil
	}
	var result interface{}
	e = c.opts.coder.Decode(b, &result)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("GetBean(%s,%s) error: %s\n", tableName, id, e)
		}
		return nil
	}
	return result
}
func (c *Cache) PutBean(tableName string, id string, obj interface{}) {
	value, e := c.opts.coder.Encode(obj)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("PutBean(%s,%s,%v) error: %s\n", tableName, id, obj, e)
		}
		return
	}
	key := c.beanKey(tableName, id)
	e = c.store.Put(key, value)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("PutBean(%s,%s,%v) error: %s\n", tableName, id, obj, e)
		}
		return
	}
}
func (c *Cache) DelBean(tableName string, id string) {
	key := c.beanKey(tableName, id)
	e := c.store.Del(key)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("DelBean(%s,%s) error: %s\n", tableName, id, e)
		}
		return
	}
}

func (c *Cache) ClearIds(tableName string) {
	prefix := c.sqlPrefix(tableName)
	e := c.store.DelPrefix(prefix)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("ClearIds(%s) error: %s\n", tableName, e)
		}
		return
	}
}
func (c *Cache) ClearBeans(tableName string) {
	prefix := c.beanPrefix(tableName)
	e := c.store.DelPrefix(prefix)
	if e != nil {
		if c.opts.logger != nil {
			c.opts.logger.Printf("ClearBeans(%s) error: %s\n", tableName, e)
		}
		return
	}
}
