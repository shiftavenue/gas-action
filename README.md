# Google Apps Script GitHub Action

> Apps Script 🤝 GitHub

This actions enables the management of [Apps Script](https://www.google.com/script/start/) code in GitHub. It can be used to do the following:

- push the code to Apps Script projects
- deploy a script
- trigger execution of a script (**NOTE:** only works if it is deployed as an API executable)

## Prerequisites

Before using this action, there are a few steps required to be able to use it. Without fulfilling all prerequisites, the action will not work. So here we go:

1. Create a [GCP project](https://cloud.google.com/resource-manager/docs/creating-managing-projects?hl=en)
1. The account that will be later used for the impersonation (see [Usage](#usage)) needs at least the IAM role `roles/oauthconfig.editor` assigned on that project
1. Create a service account in the GCP project
1. Configure [Workload Identity Federation](https://cloud.google.com/blog/products/identity-security/enabling-keyless-authentication-from-github-actions?hl=en) for the GitHub repository using this action
1. Enable the Apps Script API (`script.googleapis.com`) in the GCP project
1. Create an [OAuth Consent Screen](https://developers.google.com/apps-script/guides/cloud-platform-projects#complete_the_oauth_consent_screen) in your GCP project. The easiest way is to configure a consent screen of the type `Internal`.
1. Configure [domain-wide delegation](https://developers.google.com/cloud-search/docs/guides/delegation) for the created service account so it can impersonate any organization account (most importantly the one that will be used for the impersonation)
1. When configuring the domain-wide delegation for the service account, set all necessary OAuth access scopes that are needed for executing your scripts
1. Create the Apps Script project using the account that will be impersonated later using the domain-wide delegation (e.g. your personal account)
1. Enable the Apps Script API in your Apps Script projects by navigating to `https://script.google.com/home/usersettings`
1. Set the GCP project as "backend" project for your Apps Script project in `https://script.google.com/home/projects/<script_id>/settings`
1. Finally, use the Apps Script project ID, the GCP project ID, the service account mail, the Workload Identity pool and provider, the user to be impersonated and all necessary access scopes in this action (see more in [Usage](#usage))

<!-- As you can see, this requires a lot of setup beforehand, however this will lead to a robust setup for your Apps Script projects. For more detailed instructions and more detailed information you can visit an article we published: TODO insert link to LinkedIn article -->

## Usage

The action is designed to be used together with [`google-github-actions/auth`](https://github.com/google-github-actions/auth) to ensure a fully automated process. The Apps Script API that is used in the background to perform the deployments and executions requires an OAuth 2.0 access token with appropriate scopes (depending on what kind of Google Workspace services you work with in your scripts). The example below will use [Workload Identity Federation (WIF)](https://cloud.google.com/blog/products/identity-security/enabling-keyless-authentication-from-github-actions?hl=en), so make sure to implement the proper WIF setup first. You don't HAVE to use WIF, you can also submit a GCP service account credential file or even manually submit the access token (somehow). Using the awesome Google action in conjunction with WIF is certainly recommended though. So here we go:

```yaml
# Example for pushing script code
jobs:
  push-script-code:
    runs-on: ubuntu-latest
    permissions:
      contents: "read"
      # Is necessary to create access tokens
      id-token: "write"
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      # Authenticate to GCP using your custom project and service account
      - uses: google-github-actions/auth@v2
        id: google-auth
        with:
          project_id: my-gcp-project
          workload_identity_provider: projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider
          service_account: my-svc-account@my-gcp-project.iam.gserviceaccount.com
          # User to be impersonated via domain-wide-delegation, take user for which the Apps Script projects have been created
          access_token_subject: my-user@example.com
          token_format: access_token
          # These access token scopes are mandatory when you want to update the script content
          # see https://developers.google.com/apps-script/api/reference/rest/v1/projects/updateContent 
          access_token_scopes: "https://www.googleapis.com/auth/script.projects"

      - name: Push script code
        uses: shiftavenue/gas-action@v0
        with:
          # Use push command
          command: push
          project-id: my-apps-script-project
          access-token: ${{ steps.google-auth.outputs.access_token }}
          script-dir: ./script
```

For more detailed examples for each of the available commands (push, deploy, run), you can check out the [`examples`](./examples/) directory.

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|:-------:|
| `command` | Specifies the operation that should be performed (e.g. push code, deploy or run script). Currently allowed values are: push, deploy, run | yes | - |
| `access-token` | OAuth2 access token associated with the Google Cloud account/user used for this action | yes | - |
| `project-id` | Apps Script project ID in which the script is located | yes | - |
| `script-dir` | Directory that contains the script code | no | `"./"` |
| `function` | Entrypoint function of the Apps Script script used when triggering an execution | no | `""` |

## Outputs

This action has no outputs defined.

## Contributing

Contributions are very welcome! If you want to improve this action or if there are any problems when using it, feel free to open an [issue](https://github.com/shiftavenue/gas-action/issues) or a [pull request](https://github.com/shiftavenue/gas-action/pulls).
