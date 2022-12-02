// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0
import { Construct } from "constructs";
import { App, TerraformStack, CloudBackend, NamedCloudWorkspace, TerraformAsset, AssetType } from "cdktf";
import * as google from '@cdktf/provider-google';
import * as path from 'path';

const project = 'animated-fortnight';
const region = 'asia-northeast1';
const repository = 'animated-fortnight';

class MyStack extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new google.provider.GoogleProvider(this, 'google', {
      project,
      region,
    });

    const my_bucket = new google.storageBucket.StorageBucket(this, 'my_bucket', {
      location: region,
      name: `my-bucket-${project}`,
    });

    const my_service_account = new google.serviceAccount.ServiceAccount(this, 'my_service_account', {
      accountId: 'my-service-account',
      displayName: 'service account for this application',
    });

    for(const service of ['randomgen', 'multiply']){

      const my_asset = new TerraformAsset(this, `${service}_asset`, {
        path: path.resolve(service),
        type: AssetType.ARCHIVE,
      });
  
      const my_object = new google.storageBucketObject.StorageBucketObject(this, `${service}_object`, {
        bucket: my_bucket.name,
        name: `${my_asset.assetHash}.zip`,
        source: my_asset.path,
      });
  
      const my_function = new google.cloudfunctions2Function.Cloudfunctions2Function(this, `${service}_function`, {
        buildConfig: {
          entryPoint: service,
          runtime: 'go119',
          source: {
            storageSource: {
              bucket: my_bucket.name,
              object: my_object.name,
            },
          },
        },
        location: region,
        name: service,
        serviceConfig: {
          minInstanceCount: 0,
          maxInstanceCount: 1,
          serviceAccountEmail: my_service_account.email,
        },
      });
  
      new google.cloudfunctions2FunctionIamBinding.Cloudfunctions2FunctionIamBinding(this, `${service}_functions_iam_binding`, {
        cloudFunction: my_function.name,
        location: region,
        members: [`serviceAccount:${my_service_account.email}`],
        role: 'roles/cloudfunctions.invoker',
      });

      new google.cloudRunServiceIamBinding.CloudRunServiceIamBinding(this, `${service}_run_iam_binding`, {
        location: region,
        members: [`serviceAccount:${my_service_account.email}`],
        role: 'roles/run.invoker',
        service: my_function.name,
      });

    }

    new google.cloudbuildTrigger.CloudbuildTrigger(this, 'my_build_trigger', {
      filename: 'cloudbuild.yaml',
      github: {
        owner: 'hsmtkk',
        name: repository,
        push: {
          branch: 'main',
        },
      },
    });

  }
}

const app = new App();
const stack = new MyStack(app, "animated-fortnight");
new CloudBackend(stack, {
  hostname: "app.terraform.io",
  organization: "hsmtkkdefault",
  workspaces: new NamedCloudWorkspace("animated-fortnight")
});
app.synth();

/*

    // Cloud Run permission must be provided

    const data_policy = new google.dataGoogleIamPolicy.DataGoogleIamPolicy(this, 'data_policy', {
      binding: [{
        role: 'roles/run.invoker',
        members: ['allUsers'],
      }],
    });

    new google.cloudRunServiceIamPolicy.CloudRunServiceIamPolicy(this, 'my_policy', {
      service: my_function.name,
      location: region,
      policyData: data_policy.policyData,
    });

  }
}

const app = new App();
const stack = new MyStack(app, "fluffy-carnival");
new CloudBackend(stack, {
  hostname: "app.terraform.io",
  organization: "hsmtkkdefault",
  workspaces: new NamedCloudWorkspace("fluffy-carnival")
});
app.synth();


*/