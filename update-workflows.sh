#!/bin/sh
apt-get -y install jq
randomgen_url=`gcloud functions describe randomgen --format json --gen2 --region asia-northeast1 | jq .serviceConfig.uri`
multiply_url=`gcloud functions describe multiply --format json --gen2 --region asia-northeast1 | jq .serviceConfig.uri`
sed -i 's|var_randomgen_url|${randomgen_url}|' workflows.yaml
sed -i 's|var_multiply_url|${multiply_url}|' workflows.yaml
cat workflows.yaml
gcloud workflows deploy my-workflows --region asia-northeast1 --source workflows.yaml
