from contextlib import closing, contextmanager

try:
    from mysql import connector as mysql_connector
except Exception:
    mysql_connector = None

try:
    import psycopg2.pool
except Exception:
    psycopg2 = None

from demo import config

MYSQL_POOL = None
PGSQL_POOL = None


@contextmanager
def pgsql_pool_returner(cxn):
    try:
        yield cxn
    finally:
        PGSQL_POOL.putconn(cxn)


def mysql():
    return closing(mysql_connector.MySQLConnection(**config.MYSQL_DSN))


def pgsql():
    # return closing(psycopg2.connect(**config.PGSQL_DSN))
    global PGSQL_POOL

    if PGSQL_POOL is None:
        PGSQL_POOL = psycopg2.pool.SimpleConnectionPool(5, 10, **config.PGSQL_DSN)

    return pgsql_pool_returner(PGSQL_POOL.getconn())
