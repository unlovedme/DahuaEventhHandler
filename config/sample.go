package config

var yannotated = `# Handlers know how to send notifications to specific services.
handler:
  slack:
    # Slack "legacy" API token.
    token: ""
    # Slack channel.
    channel: ""
    # Title of the messag