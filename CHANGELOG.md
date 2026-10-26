# Changelog

All notable changes to the Pulumi provider for Forumulate. The provider is bridged from
[terraform-provider-forumulate](https://github.com/forumulate/terraform-provider-forumulate), so
every Terraform provider change appears here after `make tfgen build_sdks`.

## Unreleased

- First version: the status page resources (`forumulate:index:StatusComponent`,
  `StatusGroups`, `StatusIncidentTemplate`, `StatusMetric`) and the `getStatusPage` function, with
  SDKs for Node.js (`@forumulate/pulumi`), Python (`pulumi_forumulate`), Go and .NET.
