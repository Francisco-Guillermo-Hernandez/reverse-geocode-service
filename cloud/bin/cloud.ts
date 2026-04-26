#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib/core';
import { ReverseGeocoderServiceStack } from '../lib/reverse-geocoder-service-stack';

const app = new cdk.App();
new ReverseGeocoderServiceStack(app, 'ReverseGeocodeServiceStack');
