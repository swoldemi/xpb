# Cross-Project Billing (Google Cloud Platform)
Tired of paying Google to be your cloud provider? Want to use your mom's credit card without being charged for resources? Has your startup accelerator stopped giving you promo codes? Tired of migrating 50 Kubernetes Clusters at the end of your $300 USD Google Cloud Platform trial quota? Have a side-hustle that needs free infrastructure? Well look no further. XPB is for you!

## Expectations
- You have a GCP project with trial billing enabled (this requires a credit card)
- You are coming to the end of your $300 trial quota (not really necessary, but would defeat the purpose of the project)
- You don't want to [Upgrade](https://www.youtube.com/watch?v=OmXuQxQ-4EE) your account to charge your linked credit or debit card (yay free stuff)

## Resource Requirements
1. A GCP account, on trial, with billing enabled
2. A second GCP account, on trial, with billing enabled.
    - You may use the same billing card that was used in the first account
    - You may simply make a new Gmail account
3. Gmail access to the root accounts bound to both GCP accounts
    - They must be different accounts
    - Projects should be under different organizations
    - They must be Gmail accounts (other domains/Gsuite **not tested** -- code or test contributions are welcome!)

## Who does this benefit?
1. Startups rejected from the [Google Cloud for Startups program](https://cloud.google.com/developers/startups/) looking to lower their net burn rate
2. Students experimenting a little too much
3. Curious cloud developers

## Usage & Config
- This tool supports 2 ways of authentication
1. Application Default Credentials
2. Service account credentials

## Notes & Disclaimers 
- This project is Google Cloud Platform specific

- If I've offended anyone at Google or broken any ToS regulations, feel free to [open an issue](https://github.com/swoldemi/xpb/issues) and send that cease and desist notice my way.

- You will still need to migrate accounts (or upgrade your account) at the end of the 365 day trial. Otherwise, all resources will be shutdown.

- If you mine crypto on cloud infrastructure, __I look down on you__ (~~unless it's Oracle's cloud~~)

- Companies, however small, that have profitable software should have no reason to use this

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
    - I did it manually though the GCP Console and wanted to automated it
    - I had an excuse to use Go  

gcloud auth application-default login # Select `host` account
export GOOGLE_APPLICATION_CREDENTIALS=$(gcloud auth application-default print-access-token)
gcloud config set project YOUR_PROJECT_NAME