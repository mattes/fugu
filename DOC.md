fugu
====

## fugu.yml

Create a ``fugu.yml`` file next to your Dockerfile.

```
image: name # always mandatory
[docker-options]
```

You can use labels for different settings. 
It is possible to inherit options from other labels. 
See [advanced example](fugu.example.yml).

```
[label]
  image: name # always mandatory
  [docker-options]
```

Both the fugufile and the cli support both short and long option names. Try to
use the long option names (``env`` instead of ``e``) in the fugufile though.
Docker options that can be set multiple times (i.e. ``env``, ``link``) must
be set as array in the fugufile. See example below.


Example fugu.yml

```
image:  mattes/hello-world-nginx # mandatory
name:   hello-world-nginx
detach: true
publish:
  - 8080:80
  - 8090:81
```

## fugu run

this executes docker run

```
fugu run [fugufile] [label] [docker-options] [command] [args...]

 * fugufile defaults to fugu.yml, fugu.yaml, .fugu.yml or .fugu.yaml
 * label defaults to "default". if there is no "default" label in fugufile, use the first label found.
 * all docker run options are supported (i.e. --attach or -a)
   see http://docs.docker.com/reference/commandline/cli/#run
 * you can use `--image` to specify the docker image to be used via the cli
```


## fugu build

this executes docker build

```
fugu build [fugufile] [label] [docker-options]

 * fugufile, see above
 * label, see above
 * all docker build options are supported (i.e. --no-cache)
   see http://docs.docker.com/reference/commandline/cli/#build
 * `--tag` option (if not set via cli) will be filled with `image` variable from fugufile
 * you can use `--path` to build the image from source code at path. defaults to current directoy.
```
