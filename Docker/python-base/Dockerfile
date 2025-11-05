FROM ubuntu:22.04

MAINTAINER Jeff Morgan (jeff.morgan@appdynamics.com)

################## BEGIN INSTALLATION ######################

# Add appdynamics user
RUN useradd --create-home appdynamics
RUN mkdir /appd
#RUN chown -R appdynamics.appdynamics /appdynamics

# Deploy libpam fix (allows cron to work)
RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get -y build-dep pam
RUN export CONFIGURE_OPTS=--disable-audit && cd /root && apt-get -b source pam && dpkg -i libpam-doc*.deb libpam-modules*.deb libpam-runtime*.deb libpam0g*.deb

# Install Python 3.11 and packages
RUN apt-get update && apt-get install -y \
    python3.11 \
    python3.11-dev \
    python3-pip \
    postgresql \
    postgresql-client \
    libpq-dev \
    mysql-server \
    python3-mysqldb \
    siege \
    build-essential \
    git \
    wget \
    zip \
    unzip \
    nmap \
    locate

RUN updatedb

# pip and setuptools come with Python 3, no need for distribute
RUN pip3 install virtualenv

##################### INSTALLATION END #####################
