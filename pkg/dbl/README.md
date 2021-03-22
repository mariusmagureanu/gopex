## Database layer (DBL) package

The dbl package is responsible of persistence handling. It is based on [gorm](https://gorm.io) as on orm and
used the types defined in the ``ds`` package as model for persistence.

Currently, the dbl supports two dialects: ``postgres`` and ``sqlite``.