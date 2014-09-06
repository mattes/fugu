fugu
====

__fugu helps you to quickly run a container by storing arguments in a YAML file.__

__Why?__ We are working on [developermail.io](https://developermail.io) atm. 
The project uses a microservice architecture and consists of >15 docker images. 
During development a docker container is built, run and destroyed quite often.
With fugu we can speed up this workflow, because all ``docker run`` arguments
are stored in the ``fugu.yml`` file. We also used to put ``docker run`` statements 
in README.md, but the format wasn't consistent. Now ``fugu.yml`` is our second point of contact 
(after the Dockerfile itself), when looking at a new docker image. 
We didn't want to use fig, because the set of containers we run during
development changes often and we didn't want to keep one fig.yml for every
possible docker container combination.

Please note: fugu is not an orchestration tool like [fig](https://github.com/docker/fig). 
It is just a simple wrapper around ``docker run``.


# Installation

```bash
go get github.com/mattes/fugu

# TODO: release pre-compiled versions
```


# Usage

1) Create a ``fugu.yml`` file (maybe next to Dockerfile) and specify ``docker run``
([options](http://docs.docker.com/reference/commandline/cli/#run)). 
Valid variables are ``image``, ``command``, ``args``, and all other option variables
like ``publish`` or ``name``. The YAML file looks nicer if you don't use the
one-letter alias variables.

```yml
name: hello-world-nginx
image: mattes/hello-world-nginx
detach: true
publish: 
  - 80:80
```

2) Use ``fugu`` to run container

```bash
# in directory where fugu.yml is saved
fugu run
```

3) Profit!


# Advanced usage

```bash
fugu run [fugu.yml-path] [label] [docker-run-options] [image] [command] [args]
fugu build [fugu.yml-path] [label] [docker-build-options] [path=pwd|url|-]
```

Labels can be used for different configuration settings. See an example ``fugu.yml``:

```yml
production: &production
  name: hello-world-nginx
  image: mattes/hello-world-nginx
  detach: true
  publish: 
    - 80:80

development:
  <<: *production # inherit from production label
  detach: false
  tty: true
  interactive: true
  publish:
    - 8080:80
```

When no label argument is given, fugu will use the first label found in ``fugu.yml``.
