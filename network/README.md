## Contractor Network Config

This is an basic network for Contractor app.
There are multiple organizations:
* Orderer organization - FOI
* Development organization - Awesome Agency
* Contractor organization - Pharmatic

#### Instructions

* Run `generate.sh`
  * Generates cryptographic material
  * Creates genesis blocks for orderer and other organizations
  * Creates `docker-compose.yaml` from `docker-compose-template.yaml`

* Run `start.sh`
  * Removes any previously created networks
  * Starts network by running `docker-compose`
  * Create and join organization to the *default* channel

* Run ``stop.sh``
  * Temporarily stop network

* Run ``teardown.sh``.
  * Removes any previous networks
