# Changelog

## v0.9.0 (15 Jul 2022)

 - Feature: Introduce `kubeconfig` command for configuring kubectl access to cluster
 - Change: Improve readability of booleans in table view
 - Change: Improve readability of nested data on commands that override table view
 - Change: Revamp waiters error handling
 - Change: Revamp ssh command error handling
 - Change: Improve help messages for tag/untag commands
 - Bugfix: Display only ssh-able machines on ssh's autocomplete
 - Bugfix: Fix missing entries on autocomplete

## v0.8.0 (24 Jun 2022)

 - Feature: Add tag/untag support for resources
 - Feature: Fetch the raw credentials for clusters, keys and secrets with a CLI argument
 - Feature: Add delete-context & rename-context
 - Feature: Add secrets support (Mist v5)
 - Feature: Add `at` parameter for fetching resources
 - Feature: Add metering for volumes
 - Change: Set newly created context as the default
 - Bugfix: Fix meter command

## v0.7.5 (13 Apr 2022)

 - Bugfix: Fix adding initial context
 - Bugfix: Fix missing API doc URL
 
## v0.7.4 (24 Mar 2022)

 - Bugfix: Fix machine actions
 - Bugfix: Display all available examples in help
 
## v0.7.3 (23 Mar 2022)

 - Bugfix: Restore resource aliases
 - Bugfix: Add missing fields
 - Change: Improve help messages

## v0.7.2 (20 Dec 2021)

 - Change: Improve help messages
 - Change: Add delete_domain_image query param to undefine action
 - Bugfix: Trim always trailing slash on server URLs

## v0.7.1 (4 Dec 2021)

 - Bugfix: Fix meter command

## v0.7.0 (29 Nov 2021)

 - Change: Use spaces instead of dashes on all commands

## v0.6.0 (22 Nov 2021)

 - Feature: Introduce version command

## v0.5.0 (22 Oct 2021)

 - Change: Update create-machine params

## v0.4.0 (04 Aug 2021)

 - Feature: Add support for Kubernetes Clusters

## v0.3.0 (21 Jul 2021)

 - Feature: Introduce `meter` command for metering resources
 - Feature: Get datapoints
 - Feature: Support users & orgs endpoints
 - Feature: Add waiters in `create machine` command
 - Feature: Add support for K8s clouds
 - Feature: Add support for `create machine` on Linode, Azure, GCE
 - Feature: Add support for `get job` command
 - Change: Update openapi-cli-generator

## v0.2.0 (20 May 2021)

 - Feature: Add support for `create machine` on DigitalOcean, Equinix Metal, AWS, KVM
 - Change: Update openapi-cli-generator
 - Change: Drop GigG8 support

## v0.1.0 (11 Mar 2021)

 - Change: placeholder
 - Feature: Automate release pipelines
 - Feature: Generate Windows binaries

## v0.0.0 (2 Mar 2021)

Initial release
