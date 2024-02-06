# AARTI - Artifact Registry 

Aarti is an open source Self Hosted Custom Artifact Registry Manager that Stores data in any OCI Compliant Registries. Just like [NPM](https://www.npmjs.com) or [Composer](https://getcomposer.org) you can create your own Custom Artifact Registry using AARTI.

Also Aarti stores the data in any OCI Compliant Registry at the backend, like [Docker Hub](https://hub.docker.com) or [GCR](https://cloud.google.com/artifact-registry) etc. So you do not need to manage the database and storage separately. 


## Use Cases

### Registries for Custom Artifacts
You might need to create your Custom Artifact Registries if you want to store your modules or packages that follow a specific structure. Like Nodejs has NPM Modules, GoLang has Go Modules, HELM, Maven, etc have specified modules. Similarly you can have your own custom Artifacts. Like we at ArrowAI had specific Modules for our UI which are compiled build UI components that can be called by UI. 

In any case a Custom Artifact needs a custom Registry to be handled in a specific way including listing, downloading, creating etc. This is where Aarti comes in picture and is able to handle Registry for any kind of Custom Artifact by just adding a small module.

### Private Registries
Sometimes you need to use Private Registries for their internal Private modules, that they do not want to expose to the world.

### Registry Caching
Sometimes you want to have a Private Registry to handle caching of Modules. For example you might want a Local NPM Registry so that you can cache your NPM Modules to be available locally without going to NPM Registry. In that case you can use Aarti and a local OCI Docker Registry to run a full local experience.


### Credits

Aarti is inspired by [NPM](https://www.npmjs.com) and [Composer](https://getcomposer.org). We are also using the [OCI Artifacts](https://github.com/opencontainers/artifacts) project for the OCI Artifacts implementation. Also we are using the [OCI Distribution](https://github.com/opencontainers/distribution-spec) project for the OCI Registry implementation. Also we are using the Linka Cloud [OCI Artifact Registry](https://github.com/linka-cloud/artifact-registry) project for the OCI Artifact Registry implementation.
