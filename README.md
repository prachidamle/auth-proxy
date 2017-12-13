auth-proxy
========

A microservice that does micro things.

## Building

`make`


## Running

`./bin/auth-proxy --cluster-config=/Users/prachi/.kube/config --httpHost=localhost:9998`

## API support

1) User Login: 
   
   POST /v3/tokens?action=login
   
   JSON Input: github.com/rancher/types/apis/management.cattle.io/v3.LoginInput{}
   
   JSON Output: github.com/rancher/types/apis/management.cattle.io/v3.Token{}

2) User Logout: 
   
   POST /v3/tokens?action=logout
   
   Set Valid login TokenID in cookie 'rAuthnSessionToken' to delete the token

3) Create Derived Token: 
   
   POST /v3/tokens

   Set Valid login TokenID in cookie 'rAuthnSessionToken'
   
   JSON Input: github.com/rancher/types/apis/management.cattle.io/v3.Token{}
   
   JSON Output: github.com/rancher/types/apis/management.cattle.io/v3.Token{}

4) List Tokens: 
   
   GET /v3/tokens
   
   Set Valid login TokenID in cookie 'rAuthnSessionToken'

5) List Identities:
   
   GET /v3/identities
   
   Set Valid login TokenID in cookie 'rAuthnSessionToken'



## License
Copyright (c) 2014-2016 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
