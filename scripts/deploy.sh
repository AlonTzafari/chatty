#!/bin/bash
echo "$ALONTZ_DEV_KEY" > key.pem
chmod 400 key.pem
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i key.pem ./deployment/docker-compose.yaml admin@$ALONTZ_HOST:/home/admin/chat/docker-compose.yaml
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i key.pem admin@$ALONTZ_HOST 'docker compose -f /home/admin/chat/docker-compose.yaml up -d'
rm key.pem