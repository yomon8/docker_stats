SCRIPT_DIR=$(cd $(dirname $0) && pwd)
TMPDIR=${SCRIPT_DIR}/tmp
mkdir ${TMPDIR}
export MUNIN_PLUGSTATE=${TMPDIR}
name=docker_stats_test
docker run --name ${name} -d redis:3.2.11-alpine 
./docker_stats config
./docker_stats 
docker rm -f ${name}


