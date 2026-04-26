import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as path from "path";

export class ReverseGeocoderServiceStack extends cdk.Stack {
   constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Define the Go Lambda function
    const goLambda = new lambda.Function(this, "ReverseGeocoderServiceFunction", {
      runtime: lambda.Runtime.PROVIDED_AL2023, // Use the custom runtime
      handler: "bootstrap", // Go binary name
      code: lambda.Code.fromAsset(path.join(__dirname, "../../build")), 
      architecture: lambda.Architecture.ARM_64,
      environment: {}
    });

    // Define the Lambda function URL
    const myFunctionUrl = goLambda.addFunctionUrl({
      authType: lambda.FunctionUrlAuthType.NONE,
    });

    // Output the function's url
    new cdk.CfnOutput(this, "FunctionUrl", {
      value: myFunctionUrl.url,
    });
  }
}
