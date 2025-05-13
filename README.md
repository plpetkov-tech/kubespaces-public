# Kubespaces Public repository

Public assets for the kubespaces project, under [Apache License 2.0](./LICENSE)

## Pre-requisites

```bash
# List any pre-requisites here, for example:
# - Docker
# - Kubernetes
# - Helm
```

## Installation

```bash
# Provide installation steps here, for example:
# 1. Clone the repository
# git clone https://github.com/yourusername/kubespaces-public.git
# 2. Change to the project directory
# cd kubespaces-public
# 3. Install dependencies
# make install
```

## Usage

```bash
# Provide usage instructions here, for example:
# 1. Start the application
# make start
# 2. Access the application
# http://localhost:8080
```

## Testing

```bash
# Provide testing instructions here, for example:
# 1. Run the tests
# make test
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](./LICENSE) file for details.



## ToDo

- [ ] Add more details to the README
- [ ] Use [patches](https://www.vcluster.com/docs/vcluster/0.20.0/configure/vcluster-yaml/experimental/generic-sync?x1=1#patches-reference) to rewrite the HTTPRoute in the tenant vcluster to a proper format in the host cluster
- [ ] Add a demo for the tenant vcluster
- [ ] Add the HBONE Istio Ambient port (15008) to the NetworkPolicy created bu vcluster helm chart
- [ ] 