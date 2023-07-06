# devops
This code will setup a new VM with postgres, a load balancer, and a web app that can query from the postgres running on localhost.

I use it with google cloud and select the `e2-micro` VM so it's in the free tier.

You have to put it in `us-west1` (Oregon) or `us-central1` (Iowa).

Your boot disk has to be `standard persistent disk` up to 30GB.

So you get a 30GB hard drive, 1GB ram, and two `AMD EPYC 7B12 2250 MHz processors`. 
That's enough to run a nice little site with plenty of traffic. 
I've had links on the front page of hacker news and it never went down, didn't
even max out the cpu.

# balancer
The load balancer serves `one main purpose`: you can deploy a new version of the web app
with zero down time.

I don't use google's real load balancers or their real postgres because, free!

Without the balancer if you deployed the web app, for a second or two, the
reverse proxy would give the user a 500 error. That might not seem that bad but for
a production site, I want to be able to deploy many times a day and not affect
users ever.


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
