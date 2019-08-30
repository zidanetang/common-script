#!/bin/bash

#Author: Melo.tang
#start mongodb
#https://github.com/mongodb/mongo/blob/master/debian/init.d



help() {
    echo "Usage: install_mongo.sh -flags"
    echo "flags:" 
    echo "-c, --count      cluster members"
    echo "-m, --mode       cluster mode: replication, sharding, single"
    echo "-v, --version    mongodb version"
}

install() {
    sudo apt-get install -y mongodb-org=$version mongodb-org-server=$version mongodb-org-shell=$version mongodb-org-mongos=$version mongodb-org-tools=$version
}

replica() {
    id=$1
    port=`expr 27017 - $id`
    mkdir -p /etc/mongod/mongo$port
    mkdir -p /var/lib/mongodb$port
    chown -R mongodb:mongodb /var/lib/mongodb$port
    cat > /etc/mongod/mongo$port/mongod.conf << EOF
# mongod.conf

# for documentation of all options, see:
#   http://docs.mongodb.org/manual/reference/configuration-options/

# Where and how to store data.
storage:
  dbPath: /var/lib/mongodb$port
  journal:
    enabled: true

# where to write logging data.
systemLog:
  destination: file
  logAppend: true
  path: /var/log/mongodb/mongod$port.log

# network interfaces
net:
  port: $port
  bindIp: 127.0.0.1


# how the process runs
processManagement:
  timeZoneInfo: /usr/share/zoneinfo

#security:

#operationProfiling:

#replication:
replication:
   replSetName: "rs0"
#sharding:

## Enterprise-Only Options:

#auditLog:
EOF
}

init_repl() {
    cat >  replication.js << EOF
rs.initiate( {
   _id : "rs0",
   members: [
      { _id: 0, host: "127.0.0.1:27017" },
   ]
});
EOF
    for ((i=`expr $count - 1`; i>0; i--))
    do
        port=`expr 27017 - $i`
        sed -i "/127.0.0.1:27017/a \ \ \ \ \ \ { _id: $i, host: \"127.0.0.1:$port\" }," replication.js
    done
}

repo() {
    wget -qO - https://www.mongodb.org/static/pgp/server-$pgp.asc | sudo apt-key add -
    echo "deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/$pgp multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-$pgp.list
    sudo apt-get update -y
    mkdir -p /etc/mongod
}

#main

for arg in "$@"
do
    case $arg in
        -m|--mode)
        mode=$2
        shift # Remove arg
        shift # Remove value
        ;;
        -v|--version)
        version=$2
        shift
        shift
        ;;
        -c|--count)
        count=$2
        shift
        shift
        ;;
        -h|--help)
        help
        exit 0
        ;;
        #*)
        #echo "invalid option"
        #help
        #exit 127
        #;;
    esac
done
pgp=`echo $version | awk -F "." '{print$1"."$2}'`

if [ $mode = replication ]; then
    repo
    install
    for ((i=0; i<$count; i++))
    do
        replica $i
        port=`expr 27017 - $i`
        #su - mongodb -c "mongod -f /etc/mongod/mongo$port/mongod.conf &"
	DAEMON=/usr/bin/mongod
	DAEMON_OPTS="--config /etc/mongod/mongo$port/mongod.conf"
        start-stop-daemon --background --start --quiet \
                        --chuid mongodb:mongodb \
                        --exec $DAEMON -- $DAEMON_OPTS
    done
    if [ $count > 1 ]; then
        init_repl
        mongo admin replication.js
    fi
elif [ $mode = single ]; then
    repo
    install
    su - mongodb -c "mongod -f /etc/mongod.conf &"
fi
