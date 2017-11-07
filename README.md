ssh-auth
========

> privacyIDEA is an open solution for strong two-factor authentication like OTP tokens, SMS, smartphones or SSH keys. Using privacyIDEA you can enhance your existing applications like local login (PAM, Windows Credential Provider), VPN, remote access, SSH connections, access to web sites or web portals with a second factor during authentication. Thus boosting the security of your existing applications.

ssh-auth is a simple little Go program that connects to a privacyIDEA server and asks for SSH auth keys for the machine. This is used to centralize the tokens and make the user management easier.

Config options
--------------
The following config options is currently available for ssh-auth
```
./ssh-auth --help

Usage of ssh-auth:
  -hostname string
    	Hostname of server to validate (default "Andreas-MacBook-Pro.local")
  -server string
    	URL to PrivacyIDEA server. (default "http://127.0.0.1:5000")
  -user string
    	Username to validate
```

Sample configuration for CoreOS:
--------------------------------
The following config file for CoreOS will add the authenticator to CoreOS. Remember to alter the url for your privacyIDEA server.

```
#cloud-config
write_files:
  - path: /etc/ssh/sshd_config
    permissions: 0600
    owner: root:root
    content: |
      # Use most defaults for sshd configuration.
      UsePrivilegeSeparation sandbox
      Subsystem sftp internal-sftp
      UseDNS no
      PermitRootLogin no
      AllowUsers core
      PasswordAuthentication no
      ChallengeResponseAuthentication no
      #AuthorizedKeysCommand /bin/ssh-auth -server http://127.0.0.1:5000 -user %u

coreos:
  units:
    - name: docker.service
      enable: true
      command: start

runcmd:
  - /usr/bin/docker run --name ssh-auth -v /bin/ssh-auth:/bin/ssh-auth pasientskyhosting/ps-ssh-auth:latest
  - /bin/chmod 0755 /bin/ssh-auth
```
