# Cross-Project Billing (XPB) on Google Cloud Platform
Tired of paying Google to be your cloud provider? Want to use your mom's credit card without being charged for resources? Has your startup accelerator stopped giving you promo codes? Tired of migrating 50 Kubernetes Clusters at the end of your $300 USD Google Cloud Platform trial quota? Still haven't learned how to use [Terraform](https://github.com/hashicorp/terraform/) (me either)? Have a side-hustle that needs *even more* free cloud infrastructure? Well look no further. XPB is for you!

## Expectations
- You have a GCP project with trial billing enabled (this requires a credit card, debit card, or bank account linked).
- You are coming to the end of your $300 trial quota (not really necessary, but would defeat the purpose of the tool).
- You don't want to [Upgrade](https://www.youtube.com/watch?v=OmXuQxQ-4EE) your account to charge your linked credit or debit card (yay free stuff).

## Resource Requirements
1. A GCP account, on trial, with billing enabled.
2. A second GCP account, on trial, with billing enabled.
    - **Important**: You may use the same billing method that was used in the first account.

## Who does this benefit?
1. Startups that have been rejected from every cloud provider's Startup Program looking to lower their net burn rate.
2. Students experimenting a *little* too much.
3. Curious cloud developers.

<hr>
<h4 align="center">:warning: CAUTION :warning:</h4>

- If your billing account remains disabled for a protracted period, some resources might be removed from the projects associated with that account.
- If you associate an invalid or disabled billing account to the host's project, the project's resources will be deactivated.
<hr>

## Notes & Disclaimers
- This project is Google Cloud Platform specific. [Keep in mind that projects and resources under a free trial billing account are not covered by GCP's Service Level Agreements](https://cloud.google.com/terms/free-trial/) 

- If I've offended anyone at Google or broken any ToS regulations, feel free to [open an issue](https://github.com/swoldemi/xpb/issues) and send that cease and desist notice my way.

- I am not responsible for what happens to your GCP resources while using this tool, directly or otherwise.

- You will still need to migrate accounts (or upgrade your account) at the end of the 365 day trial. Otherwise, all resources will be shutdown until a valid billing account is added to the project. It's time to learn Terraform...

- If you mine crypto on cloud infrastructure, __I look down on you__ (~~unless it's Oracle's cloud~~). A majority of the time, mining software and unusual, long-term spikes in CPU usage will be detected and your project will be deleted promptly. Don't waste your time. Be ethical.

## Invocation flow
- The account running out of trial credits will be referred to as `host` (because it is the root host the project who's resources you would like to keep running). The new account with ~$300 in trial credit will be referred to as `guest`. 
[to do]

## Usage & Config
- This tool currently depends on the email and password for both the host and guest account. Additionally, it also uses Service account credentials to very some expectations of the host and gust accounts.
  - For instructions on how to obtain service account credentials for the program see: [url here]
- Enable the **Cloud Billing API** on the host and guest accounts through the GCP API Library dashboard (https://console.cloud.google.com/apis/api/cloudbilling.googleapis.com/overview)
- Enable the **Identity and Access Management** API on the host and guest accounts through the GCP API Library dashboard (https://console.cloud.google.com/apis/api/iam.googleapis.com/overview)

## FAQ
1. Does this actually work?
    - Yes!
2. Should I use Amazon Web Serivces, Google Cloud Platform, or Microsoft Azure?
    - Yes!
3. Why did you make this?
    - Why didn't you?
4. Couldn't I just use `gcloud alpha billing projects link`?
    - I dunno, can you?

## Contributing
>This is actually just a to-do list in disguise.

If you are intrested in contributing, the following would be greatly appreciated:
1. Update IAM assignments, such that, the bare-minimum privilege is granted to the guest's service account. If they are already at a state of least privledge, open a PR to edit this README!
2. Automatic enable of the Cloud Billing, Cloud Resource Manager, and Cloud IAM API's.
3. Automatic activation of the service accounts. Possibly calling `gcloud auth activate-service account --key-file={path}` using package `os/exec` or something similar.
4. Automatic authentication of both guest and host via three-legged OAuth exchange through `gcloud auth login`. Might be possible to do without a webserver, using package `os/exec`, since the gcloud CLI already runs a webserver for the redirect (see https://github.com/swoldemi/xpb/blob/master/auth.go#L84-L92).
5. Verify that *both* service account keys may be revoked without any reprocussion. Will the billing account become unlinked?
6. Better functionality of waiting for elements to load, intead of calling time.Sleep. The slower your internet, the longer the sleep will need to be.

## API Calls
If anything, you may use this project as a reference for the structure and usage of the Google's generated APIs https://github.com/googleapis/google-api-go-client. [Create an issue](https://github.com/swoldemi/xpb/issues) if you have any questions.
  

// Move this
1. On the host account, create a service account at https://console.cloud.google.com/iam-admin/serviceaccounts/create?authuser=1&project=nickel-api. Grant it the Project Billing Manager and Project IAM Admin roles. Do the same on the guest account.
[image here]

2. Activate the service accounts by executing 
```bash
$ gcloud auth activate-service-account --key-file=[HOST_KEY_FILE_PATH]
$ gcloud auth activate-service-account --key-file=[GUEST_KEY_FILE_PATH]
```

### Notes
> The fact that this project depends on `https://github.com/tebeka/selenium` which depends on `https://github.com/BurntSushi/xgb` is a complete coincidence.
> As of July 2019, the non-organizational, GCP project used during the development of XPB is currently <100 days into the GCP trial, but has consumed over $700 in credits (~2.3 billing accounts worth)
> You can also use promo codes, albeit harder to come by, to your trial accounts. Promo codes are applied to accounts as an additional, auxiliary billing account.