# devops
This code will setup a new VM with postgres, a load balancer, and a web app that can query from the postgres running on localhost.

I use it with [google cloud](https://cloud.google.com/) and select the `e2-micro` VM so it's in the free tier.

You have to put it in `us-west1` (Oregon) or `us-central1` (Iowa).

Your boot disk has to be `standard persistent disk` up to 30GB.

So you get a 30GB hard drive, 1GB ram, and two `AMD EPYC 7B12 2250 MHz processors`. 
That's enough to run a nice little site with plenty of traffic. 
I've had links on the front page of hacker news and it never went down, didn't
even max out the cpu.

# balancer
The load balancer serves `one` main purpose: you can deploy a new version of the web app
with zero down time.

I don't use google's real load balancers or their real postgres because, free!

Without the balancer if you deployed the web app, for a second or two, the
reverse proxy would give the user a 500 error. That might not seem that bad but for
a production site, I want to be able to deploy many times a day and not affect
users ever.

The other purpose is to run on port 443 and handle SSL and the certs from
[letsencrypt](https://letsencrypt.org/). It also runs on port 80 and just 
forwards any request on 80 to 443. i.e. you can't make an http request.
Everything forwards to https.

So it runs two golang [httputil.ReverseProxy](https://pkg.go.dev/net/http/httputil#NewSingleHostReverseProxy) things. One on port 3000 and another on 3001.

The balancer sends 100% of traffic to either 3000 or 3001. It never splits up
the traffic 50% to each because that's not the point. The point is to be able to
do this:

```
scp web aa@YOUR-IP:
ssh aa@YOUR-IP 'bash -s' < script-3001.sh
```

and script-3001.sh is:

```
systemctl stop web-3001.service
mv /home/aa/web /home/aa/web-3001
systemctl start web-3001.service
```

It's safe to run `systemctl stop web-3001.service` because 0% of traffic is
going it to. 100% is going to 3000. That's the default and how it starts.

The logic to do a deploy hits your site at a special url. The url has a guid
so no one will be able to guess this url and use it but you.

Like this:

```
https://many.pw/f0e3267a-376c-4a21-8f53-f4b5192357c6/3000
```

Note `f0e3267a-376c-4a21-8f53-f4b5192357c6` is not my real guid! Keep your
guid secret. Anyway that route of /guid/3000 makes the balancer change to
the other one it's not using. If it's using 3000 it goes to 3001. If it's
on 3001 that command makes it go back to 3000.

This is done with:

```
if WebPort == 3000 {
  WebPort++
} else {
  WebPort--
}
ReverseProxyWeb = makeReverseProxy(WebPort, false)
```

So there is a little script to query the current WebPort value and know
which one is safe to run `systemctl stop web-%s.service` on where %s gets
filled in as either 3000 or 3001

There are two systemd service files:

```
ExecStart=/home/aa/web-3000 run 3000
```
```
ExecStart=/home/aa/web-3001 run 3001
```

Notice they use a different binary. This allows us to scp a file called `web`
to /home/aa/ and then call `systemctl stop` on the right service. THEN
you can `mv web web-3000` or `mv web web-3001` because if the service is
running you CANNOT replace the binary.

You build on your local machine:

```
GOOS=linux GOARCH=amd64 go build
```

And boom, upload new version, stop the right service, rename the file,
start back up the new service, and then hit that special URL and all
of a sudden users get the new version!

# key gen
You need to run:

```
# save as ~/.ssh/aa
ssh-keygen -t ed25519 -C aa@devops
# save as ~/.ssh/root
ssh-keygen -t ed25519 -C root@devops
```

You could change `aa` to whatever username you like, but I say might as
well just use aa it's a nice username.

The public version `aa.pub` and `root.pub` need to be added to your google
VM's list of SSH keys.

# env vars

```
export BALANCER_GUID=?
export VM_IP=YOUR-IP
```

You can get your IP from the google VM (it's free!). Isn't that amazing. A free IP in
today's world.

Your BALANCER_GUID value you will set after running a command.

# How to run

First make sure you run `./build.sh` in the `balancer` dir and then the `web` 
directory so your binary files are ready. Then run each of commands below one by one.

The `./vm env youremail yourdomains` one is very special. Your email is the email
you want to use for the letsencrypt cert. And yourdomains is a comma separated list
of domains. For example:

```
./vm env fred@gmail.com good.com,great.co,other-domain.org
```

As long as your have an A record pointing to the IP of your google VM in
each domain,
letsencrypt will be able to make a cert.

So I like to run this list one by one but you could also place these in
a file and run it all at once!

It also picks a random guid for you and outputs:

```
fmt.Println("export BALANCER_GUID=" + guid)
```

So that's how you get that value.

```
./vm psql
./vm cp ../aa.conf /etc/systemd/system/ root
./vm env youremail yourdomains
./vm cp ../balancer/balancer.service /etc/systemd/system/ root
./vm cp ../balancer/balancer /home/aa/ aa
./vm reload balancer
./vm cp ../web/web /home/aa/web-3000 aa
./vm cp ../web/web /home/aa/web-3001 aa
./vm web
./vm reload web-3000
./vm reload web-3001
```

# deploy

```
./vm deploy-web yourdomain
```

That's all you need to run to hit that special url after it uploads the new binary.

And you should not need to do this often because it will require downtime but
if you want to make logic changes to the balancer you can deploy it with:

```
./vm deploy-balancer
```
