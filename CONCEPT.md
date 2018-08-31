# Concept

Features:
- clone all related repos from a single descriptor file
- creates it's own multi-repo docker-compose config... should be 100% usable by docker-compose itself
- switch between projects on demand
- must be cross-platform?
- seamlessly allow for repos to be a part of multiple projects


TODO:
- formatting UI output
- logging

## Reference tools

docker-compose
git
git-town
git-runner
bashful


## Example usage:

stitch init

...high level project verbs...
```bash
stitch new <git-url> [name]
stitch delete <name>
stitch list
stitch describe <name>  # shows the repos + paths of a project. Additionally shows the services of each repo
stitch switch <name>
# what about showing the current project?
```

...project service commands...
```bash
stitch build
stitch up [service...]
stitch down 
stitch start [service...]
stitch stop [service...]
stitch exec <service> <command>
stitch update [service...]    # git sync current branches? or just pull the latest?
stitch cd <service>
stitch bash <service>
stitch open <service> [command]
stitch list services             # shows services exposed by each repo in the current project

stitch run <script-name> [options]

stitch logs [service]
stitch ps
stitch watch   # ps in a loop
```

...maybe pipeline commands...? requires building another tool
```bash
stitch new pipeline <service>
stitch run-pipeline <service>
```

## Behind the scenes

```bash
~/.stitch
    stitch.yaml
    workspace.gob    # state file with list of projects, destinations, etc
```

## Creating a stitch project

In a standalone repo:

```yaml
# stitch-project.yaml or .stitch-project.yaml
version: 1
project: a-project-name

# invocable via 'run'
scripts:
  clean: 'rm -rf something??' # how do we know which project to run this against, or which path we're in?...
  # this would solve that problem, optional nesting
  clean2:
  - cmd: 'rm -rf the stuff'
    repo: best-app 
  - cmd: 'rm -rf the other stuff blerg'
    repo: another-cool-app
  calico: 'echo "cats are cooool!"'

repos:
# some day make it from a gist...
- name: a gist example
  gist: https://gist.github.com/wagoodman/22d3cd9724507c3d723bca861e02c8bc#file-app

- name: best-app
  git: git@github.com:someone/the-best-app.git
  # OPTIONAL: a git-sha or git-tag to the specific version
  version: 358da8fe494f9d6ad9d67a9382
  # maybe you have all project dev on this branch...
  branch: dev-project-name

- name: another-cool-app
  git: git@github.com:someone/just-another-cool-app.git
  
  # there is an assumed path to clone the repos, unless overridden here
  path: $GOPATH/src/github.com/someone/just-another-cool-app


  # since the repo is missing the stitch.yaml file, then that can be provided here.
  # This will act as an overlay on top of an existing config to override values as well.
  stitch-config:  # or just stitch.yaml: ??
    services:
      # this is the default inclusion
      include:
        - .*
      # but we can ignore some of these services
      exclude:
        - proxy
    ...

# we can sepcify this here, or in another file altogether (or not at all).
# THIS WILL BE A PROBLEM... as some files in the mount will be out of context entirely? like files from other repos, nginx config snippets.
docker-compose:
  services:
    proxy:
      image: nginx:1.15-alpine
      networks:
        lan:
      ports:
        - 80:80
      volumes:
        - ./config/nginx.conf:/etc/nginx/nginx.conf:ro
      depends_on:
        - app
```

From `the-best-app` repo:

```
.
├── .stitch.yaml
├── dockerfiles
│   ├── Dockerfile.app
│   ├── Dockerfile.pipeline
│   └── Dockerfile.toolchain
├── docker-compose.yaml
├── [docker-compose.setup.yaml]
├── src
└── ...


```

```yaml
# docker-compose.yaml
services:
  db:
    image: postgres:9.6.2-alpine
    networks:
      lan:
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: thebestdb
      PGDATA: /pgdata
    volumes:
      - db-pgdata:/var/lib/postgresql/data

  proxy:
    image: nginx:1.15-alpine
    networks:
      lan:
    ports:
      - 80:80
    volumes:
      - ./config/nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - app

  app:
    image: the-best-app:latest
    build:
      context: .
      dockerfile: dockerfiles/Dockerfile.app
      args:
        USER_UID: ${USER_UID}
    networks:
      lan:
    command: yarn start
    working_dir: /app
    ports:
      - "3000:3000"
      - "4444:4444"
      - "9222:9222"
    volumes:
      - .:/app
      - yarn_cache:/home/node/.cache
      - node_modules:/app/node_modules
    depends_on:
      - db
```


Why would we want a stitch file here? as a slot of pluggable options? but this means that if I want a repo in a project it needs a config... not desirable.
HOWEVER! we can allow hinting here in a file... or also allow this hinting in the stitch project file instead (so this repo doens't **need** any modifications)
So... this file is entirely optional.
```yaml
# stitch.yaml or .stitch.yaml
version: 1
  services:
    # this is the default inclusion
    include:
      - .*
    # services in the compose file that aren't meant to be used when stitching
    exclude:
      - proxy
```



## Open questions

- how will cross-repo proxy configuration work? Can it be automagical? ... lets assume nginx + config snippets for now since there is a default.d dir
- port remapping for conflicts will not be easy since other configs/apps/services may depend on a particular value.
- nesting of projects???