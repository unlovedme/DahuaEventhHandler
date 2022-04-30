package config

var yannotated = `# Handlers know how to send notifications to specific services.
handler:
  slack:
    # Slack "legacy" API token.
    token: ""
    # Slack channel.
    channel: ""
    # Title of the message.
    title: ""
  hipchat:
    # Hipchat token.
    token: ""
    # Room name.
    room: ""
    # URL of the hipchat server.
    url: ""
  mattermost:
    room: ""
    url: ""
    username: ""
  flock:
    # URL of the flock API.
    url: ""
  webhook:
    # Webhook URL.
    url: ""
  msteams:
    # MSTeams API Webhook URL.
    webhookurl: ""
  smtp:
    # Destination e-mail address.
    to: ""
    # Sender e-mail address .
    from: ""
    # Smarthost, aka "SMTP server"; address of server used to send email.
    smarthost: ""
    # Subject of the outgoing emails.
    subject: ""
    # Extra e-mail headers to be added to all outgoing