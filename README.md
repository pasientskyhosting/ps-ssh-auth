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
    	Hostname of server to validate (default "localhost")
  -login string
    	Login username to PrivacyIDEA (default "admin")
  -pass string
    	Login password to PrivacyIDEA (default "test")
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
      UsePrivilegeSeparation sandbox
      Subsystem sftp internal-sftp
      UseDNS no
      PermitRootLogin no
      AllowUsers core
      PasswordAuthentication no
      ChallengeResponseAuthentication no
      AuthorizedKeysCommand /auth/ssh-auth -server http://192.168.5.63:5000 -login admin -pass test -user %u
      AuthorizedKeysCommandUser root

coreos:
  units:
    - name: ssh-auth.service
      command: start
      content: |     
        [Unit]
        Description=ssh-auth

        [Service]
        ExecStart=-/bin/rm -rf /auth
        ExecStart=/bin/mkdir /auth
        ExecStart=/bin/wget -O /auth/ssh-auth https://github.com/pasientskyhosting/ps-ssh-auth/releases/download/v1.0/ssh-auth
        ExecStart=/bin/chmod 0700 /auth/ssh-auth
        Type=oneshot
```
