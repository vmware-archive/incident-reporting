
export CLIENT_USER='tscanlan-dev'
export CLIENT_PASSWORD='XXXX'
export CLIENT_URL=mgmt.blockchain.vmware.com/blockchains/XXXX/api/concord/eth
export PRODUCTION_URL="https://${CLIENT_USER}:${CLIENT_PASSWORD}@${CLIENT_URL}"

docker run -e PRODUCTION_URL -it index.docker.io/tompscanlan/incident-reporting-truffle:v1.5 truffle deploy --network production --reset

# replace this with your own contract address, or use this one for an existing sample
export CLIENT_CONTRACT_ADDRESS="0xb5653804dc7e6c45c2a5e1b6fbf25d1d393f3062"
docker run -d -e CLIENT_URL -e CLIENT_USER -e CLIENT_PASSWORD -e CLIENT_CONTRACT_ADDRESS -p 8080:80 index.docker.io/tompscanlan/incident-reporting-ui:v1.5

open http://localhost:8080/logs
