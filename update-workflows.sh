#!/bin/sh
randomgen_url=`gcloud functions describe randomgen --format json --gen2 --region ${_REGION} | jq .serviceConfig.uri`
multiply_url=`gcloud functions describe multiply --format json --gen2 --region ${_REGION} | jq .serviceConfig.uri`
sed -i 's|var_randomgen_url|${randomgen_url}|' workflows.yaml
sed -i 's|var_multiply_url|${multiply_url}|' workflows.yaml
cat workflows.yaml
gcloud workflows deploy my-workflows --source workflows.yaml
