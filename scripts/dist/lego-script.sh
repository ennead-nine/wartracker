#!/bin/bash

AWS_ACCESS_KEY_ID=AKIA5VWCZQKRNW4WWCTA \
AWS_REGION="us-west-1" \
AWS_HOSTED_ZONE_ID=Z01286661IBWXA5QX34IN \
lego --email ennead.tbc@gmail.com --dns route53 -d '*.ennead.games' -d ennead.games --path ~/.certs run
