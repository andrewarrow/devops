# devops

# env vars

```
export BALANCER_GUID="dd9b7cc1-9dbf-4fcd-ab2e-fcae5d9d6a38"
export BALANCER_EMAIL="you@somewhere.com"
export BALANCER_DOMAINS="mydomain1.com,mydomain2.com"

export VM_IP="34.83.130.106"
```

# key gen
```
# save as ~/.ssh/aa
ssh-keygen -t ed25519 -C aa@devops
# save as ~/.ssh/root
ssh-keygen -t ed25519 -C root@devops
```
# Use

```
./vm psql
./vm cp ../aa.conf /etc/systemd/system/ root
./vm env
./vm cp ../balancer/balancer.service /etc/systemd/system/ root
./vm cp ../balancer/balancer /home/aa/ aa
./vm reload balancer
./vm cp ../web/web /home/aa/web-3000 aa
./vm cp ../web/web /home/aa/web-3001 aa
./vm web
./vm reload web-3000
./vm reload web-3001
```
