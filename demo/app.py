import ipaddress
import logging
import math
import os
import random
import requests
import socket
import time
import urllib.parse

from flask import Flask, render_template, render_template_string, request, abort

from demo import db

app = Flask(__name__)

log_level = os.getenv('LOG_LEVEL', 'INFO').upper()
logging.basicConfig(level=getattr(logging, log_level, logging.INFO))


class MissingArgumentException(Exception):
    status_code = 400

    def __init__(self, message):
        super(MissingArgumentException, self).__init__()
        self.message = message


@app.errorhandler(MissingArgumentException)
def handle_missing_argument_exception(error):
    return render_template("error.html", msg=error.message), error.status_code


@app.route('/')
def index():
    return render_template("index.html")


@app.route('/wave/<whatever>')
def response_time_wave(whatever):
    a = 0.2
    b = 3
    c = 1
    x = math.radians(time.time())

    try:
        c = int(whatever)
    except ValueError:
        pass

    delay = a * (math.sin(b * x + c) + 1.0)
    time.sleep(delay)
    return render_template("wave.html", delay=delay)


@app.route('/error/<when>')
def cause_error(when):
    if when == 'always' or random.randint(0, 9) == 0:
        raise random_exception()
    return render_template("ok.html")


@app.route('/query/<dbtype>')
def query_db(dbtype):
    if dbtype not in ('pgsql', 'mysql'):
        abort(400, description=f"Invalid database type: {dbtype}")

    query_type = random.choice(('slow', 'error', 'normal', 'normal', 'normal', 'normal', 'normal', 'normal', 'normal', 'normal'))
    sleep = random.randrange(1, 8) / 10.0

    if dbtype == 'pgsql':
        with db.pgsql() as cxn:
            cur = cxn.cursor()

            if query_type == 'slow':  # Be a slow query
                cur.execute("SELECT pg_sleep(%f)" % sleep)
            elif query_type == 'error':  # Be an erroneous query
                cur.execute("SELECT sql - error")
            else:
                cur.execute("SELECT 123")
    elif dbtype == 'mysql':
        with db.mysql() as cxn:
            cur = cxn.cursor()

            if query_type == 'slow':  # Be a slow query
                cur.execute("SELECT SLEEP(%f)" % sleep)
            elif query_type == 'error':  # Be an erroneous query
                cur.execute("SELECT sql - error")
            else:
                cur.execute("SELECT 123")

    return render_template("query.html", dbtype=dbtype, type=query_type)


def is_safe_url(url):
    """Validate URL to prevent SSRF attacks."""
    parsed = urllib.parse.urlparse(url)
    
    if parsed.scheme not in ('http', 'https'):
        logging.warning(f"URL validation failed: invalid scheme {parsed.scheme}")
        return False, "URL not allowed"
    
    if not parsed.hostname:
        logging.warning("URL validation failed: no hostname")
        return False, "URL not allowed"
    
    hostname_normalized = parsed.hostname.lower().rstrip('.')
    
    allowed_hosts = os.getenv('ALLOWED_OUTBOUND_HOSTS', '').split(',')
    allowed_hosts = [h.strip().lower().rstrip('.') for h in allowed_hosts if h.strip()]
    
    if not allowed_hosts:
        logging.error("URL validation failed: ALLOWED_OUTBOUND_HOSTS not configured")
        return False, "URL not allowed"
    
    if hostname_normalized not in allowed_hosts:
        logging.warning(f"URL validation failed: hostname {hostname_normalized} not in allowlist")
        return False, "URL not allowed"
    
    allowed_ports = os.getenv('ALLOWED_OUTBOUND_PORTS', '80,443').split(',')
    allowed_ports = [int(p.strip()) for p in allowed_ports if p.strip()]
    port = parsed.port or (443 if parsed.scheme == 'https' else 80)
    
    if port not in allowed_ports:
        logging.warning(f"URL validation failed: port {port} not in allowed ports")
        return False, "URL not allowed"
    
    try:
        ip_addresses = socket.getaddrinfo(parsed.hostname, None)
    except socket.gaierror as e:
        logging.warning(f"URL validation failed: cannot resolve {parsed.hostname}: {str(e)}")
        return False, "URL not allowed"
    
    for family, _, _, _, sockaddr in ip_addresses:
        ip = sockaddr[0]
        try:
            ip_obj = ipaddress.ip_address(ip)
            if not ip_obj.is_global or ip_obj.is_private or ip_obj.is_loopback or ip_obj.is_link_local or ip_obj.is_multicast or ip_obj.is_reserved or ip_obj.is_unspecified:
                logging.warning(f"URL validation failed: {parsed.hostname} resolves to non-global IP {ip}")
                return False, "URL not allowed"
        except ValueError as e:
            logging.warning(f"URL validation failed: invalid IP {ip}: {str(e)}")
            return False, "URL not allowed"
    
    return True, None


@app.route('/http')
def http_exit_call():
    url = request.args.get('url')

    if url is None:
        raise MissingArgumentException('required argument "url" is missing')

    lower_url = url.lower()
    if not lower_url.startswith('http://') and not lower_url.startswith('https://'):
        raise MissingArgumentException('required argument "url" must be a URL with protocol, like http://...')

    is_safe, error_msg = is_safe_url(url)
    if not is_safe:
        raise MissingArgumentException(error_msg)

    session = requests.Session()
    session.trust_env = False
    
    try:
        resp = session.get(
            url, 
            timeout=(3.05, 5), 
            allow_redirects=False,
            proxies={'http': None, 'https': None}
        )
        content_length = resp.headers.get('content-length', 'n/a')
        return render_template("http_exit.html", url=url, len=content_length)
    except requests.exceptions.RequestException as e:
        logging.error(f"Request failed for {url}: {str(e)}")
        raise MissingArgumentException("Request failed")
    finally:
        session.close()


def random_exception():
    if random.randint(0, 3) == 0:
        my_dict = {}
        return my_dict['typo_eror']

    if random.randint(0, 3) == 0:
        my_list = [1, 2, 3, 4]
        return my_list[10]

    if random.randint(0, 3) == 0:
        assert True is False

    return int('abc')


if __name__ == '__main__':
    debug_mode = os.getenv('FLASK_DEBUG', '0') == '1'
    host = os.getenv('FLASK_HOST', '127.0.0.1')
    
    try:
        port = int(os.getenv('FLASK_PORT', '9000'))
    except ValueError:
        logging.warning("Invalid FLASK_PORT value, using default 9000")
        port = 9000
    
    app.run(host, port, debug=debug_mode)
