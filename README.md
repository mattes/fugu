<img src="fugu.png" height="200">

# What is fugu?

 * fugu is a wrapper around docker run & build commands.
 * fugu loads arguments from a fugu.yml file and 
   merges these arguments with cli options.
 * fugu is NOT an orchestration tool like [fig](https://github.com/docker/fig). 

 
__Example__


```yml
# fugu.yml (maybe stored next to Dockerfile)
image:  mattes/hello-world-nginx # mandatory
name:   hello-world-nginx
detach: true
publish:
  - 8080:80
```

```bash
$ fugu run -e VERY=nice # runs ...
docker run --detach --name="hello-world-nginx" --env="VERY=nice" --publish="8080:80" mattes/hello-world-nginx
```

All docker options are supported. See __[full documentation here](DOC.md)__.


# Installation

```bash
# Mac OS X
curl -L https://github.com/mattes/fugu/releases/download/v0.1.6/fugu.v0.1.6.darwin.x86_64.tar.gz | tar xvz
mv fugu.v0.1.6.darwin.x86_64 /usr/local/bin/fugu
chmod +x /usr/local/bin/fugu

# Linux
curl -L https://github.com/mattes/fugu/releases/download/v0.1.6/fugu.v0.1.6.linux.x86_64.tar.gz | tar xvz
mv fugu.v0.1.6.linux.x86_64 /usr/local/bin/fugu
chmod +x /usr/local/bin/fugu

# from source, if you have GO installed
go get github.com/mattes/fugu
```


# Why fugu?

We are working on [developermail.io](https://developermail.io) atm. 
The project uses a microservice architecture and consists of lots of docker images. 
During development a docker container is built, run and destroyed quite often.
With fugu we can speed up this workflow, because all ``docker run`` arguments
are stored in the ``fugu.yml`` file. We also used to put ``docker run`` statements 
in ``README.md``, but the format wasn't consistent. Now ``fugu.yml`` is our second point of contact 
(after the ``Dockerfile`` itself), when looking at a new docker image somebody else created.
We didn't want to use fig, because the set of containers we run during
development changes often and we didn't want to have one ``fig.yml`` for every
possible docker container combination.


---

[![Build Status](https://travis-ci.org/mattes/fugu.svg?branch=master)](https://travis-ci.org/mattes/fugu)
[![GoDoc](https://godoc.org/github.com/mattes/fugu?status.svg)](https://godoc.org/github.com/mattes/fugu)


### Credits

Thanks to [Thiago Lifter](http://www.thiagolifter.com.br) for his nice fugu fish logo.
