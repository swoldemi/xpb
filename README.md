# Cross-Project Billing (XPB) on Google Cloud Platform
Tired of paying Google to be your cloud provider? Want to use your mom's credit card without being charged for resources? Has your startup accelerator stopped giving you promo codes? Tired of migrating 50 Kubernetes Clusters at the end of your $300 USD Google Cloud Platform trial quota? Have a side-hustle that needs free infrastructure? Well look no further. XPB is for you!

## Expectations
- You have a GCP project with trial billing enabled (this requires a credit or debit card).
- You are coming to the end of your $300 trial quota (not really necessary, but would defeat the purpose of the project).
- You don't want to [Upgrade](https://www.youtube.com/watch?v=OmXuQxQ-4EE) your account to charge your linked credit or debit card (yay free stuff).

## Resource Requirements
1. A GCP account, on trial, with billing enabled
2. A second GCP account, on trial, with billing enabled.
    - Although You may use the same billing card that was used in the first account.

## Who does this benefit?
1. Startups rejected from the [Google Cloud for Startups program](https://cloud.google.com/developers/startups/) looking to lower their net burn rate.
2. Students experimenting a *little* too much.
3. Curious cloud developers.
***CAUTION***
- If your billing account remains disabled for a protracted period, some resources might be removed from the projects associated with that account.
- If you associate an invalid or disabled billing account to the host's project, it's resources will be deactivated.

## Usage & Config
- This tool currently only supports Service account credentials for authentication
  - For instructions on how to obtain service account credentials for the program see: [url here]
- Enable the Cloud Billing API on the host and guest accounts through the GCP API Library dashboard (https://console.cloud.google.com/apis/api/cloudbilling.googleapis.com/overview)
- Enable the Identity and Access Management API on the host and guest accounts through the GCP API Library dashboard (https://console.cloud.google.com/apis/api/iam.googleapis.com/overview)


## Notes & Disclaimers 
- This project is Google Cloud Platform specific. [Keep in mind that projects and resources under a free trial billing account are not covered by GCP's Service Level Agreements](https://cloud.google.com/terms/free-trial/) 

- If I've offended anyone at Google or broken any ToS regulations, feel free to [open an issue](https://github.com/swoldemi/xpb/issues) and send that cease and desist notice my way.

- I am not responsible for what happens to your GCP resource while using this tool, directly or otherwise.

- You will still need to migrate accounts (or upgrade your account) at the end of the 365 day trial. Otherwise, all resources will be shutdown until a valid billing account is added to the project.

- If you mine crypto on cloud infrastructure, __I look down on you__ (~~unless it's Oracle's cloud~~). Be ethical.

## Invocation flow
- The account running out of trial credits will be referred to as `host` (because it is the root host the project who's resources you would like to keep running). The new account with $300 in trial credit will be referred to as `guest`. 
1. Generate Application Default Credentials for the Google OAuth client
1. Authenticate with the V1 IAM API on the `host` account.
2. Invite the `guest` account to the `host`'s project and grant the `guest` the `Billing Administrator` IAM role.
3. 
4. 

## FAQ
1. Does this actually work?
    - Yes!
2. Why did you make this?
    - Why didn't you?
3. Should I use Amazon Web Serivces, Google Cloud Platform, or Microsoft Azure?
    - Yes!
4. Couldn't I just use `gcloud alpha billing projects link`?
    - Can you?

// Move this
1. On the host account, create a service account at https://console.cloud.google.com/iam-admin/serviceaccounts/create?authuser=1&project=nickel-api. Grant it the Project Billing Manager and Project IAM Admin roles. Do the same on the guest account.
[image here]

2. Activate the service accounts by executing 
```bash
$ gcloud auth activate-service-account --key-file=[HOST_KEY_FILE_PATH]
$ gcloud auth activate-service-account --key-file=[GUEST_KEY_FILE_PATH]
```