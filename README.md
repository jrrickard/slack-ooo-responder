# An out of office responder for Slack

Slack allows you to add informative status messages and snooze notifications when you are out of the office. I noticed that I was still getting a bunch of messages though, so I wrote an out of office responder!

## How to use

First, grab a slack token for your account. I think they call these [legacy tokens](https://api.slack.com/custom-integrations/legacy-tokens). 

Next, build the application or grab the docker container: gcr.io/jrrickard-178216/slack-responder:latest. There is also a sample Kubernetes manifest (and accompanying secret template).

You can provide the slack token in a config file or via an environment variable. A sample config file is included in the repo.

To run: ./slack-ooo-responder -config=<someurl> or docker run -d gcr.io/jrrickard-178216/slack-responder:latest -config <some-url>

Note: currently, timezones aren't really handled well. Define the time ranges in the config file in UTC or mount timezone/local time into the container. 


  
