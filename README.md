# NDM - Not Docker on Mac

Are you tired of a lot of different solutions to replace Docker on Mac? Yeah me too.

So I have decided to create another one and now we have 11 different solutions...

Kidding, was just an attempt to do a fast bootstrap of a Linux Machine to run Docker on Mac.

## What happens when you run this program?
In the end? You get something like this:
```
Docker Machine Started
* Please execute the following to add the SSH Key to the well-known hosts:
ssh ndm@172.16.220.133 exit

* Now just use the docker daemon there:
export DOCKER_HOST=ssh://ndm@172.16.220.133
```

And then you can use docker normally

## What works?
* Building images
* Running images
* Docker buildx should probably work as well
* As this is NAT'ed VM, if you use stuff like VPNs this internal machine will be able to do internal stuff (not sure if DNS works yet)

## What does not work yet?
So before you start running stuff and getting frustrated, let me tell you what does not work:
* Port forwarding -> You are running on another machine! That said, when you do `docker run -p ...` you should point to that machine IP
* KinD -> KinD works actually, but as KinD thinks it's running "locally" you probably need to do some ssh port forwarding, etc

## Pre Reqs:
* [Vmware Fusion](https://www.vmware.com/products/fusion.html) - Used to create the VM that will run our docker containers.

* [OVFTool](https://developer.vmware.com/tool/ovf) - Used to convert the PhotonOS OVA into a Virtual Machine

* Docker - THE CLI ONLY!! - `brew install docker` should be enought

## How this works?

This program runs some basic steps to bootstrap a VM that will run our docker engine in ALMOST a transparent way (lies...)
* Download a PhotonOS OVA (if it does not exists) into $HOME/.cache/ndm/photon.ova
* Generates some basic cloud-init to add into this machine when being created (add user, enable docker, etc)
  * Here, it uses the .ssh/id_rsa.pub key for the passwordless login. 
* Create a VM with the above cloud-init
* Starts the machine

## Running
Right now, just do a `go run .`

I plan to, if this becomes more annoying to me and other people improve this code (like unhardcoding stuff, generating binaries, adding flags, etc)

## TODO
* Generate SSH Key automatically if it does not exists
* Add the Host SSH Key automatically...maybe...on our host