# Development Secrets

This directory is for sample app development secrets. The checked-in
`development.env` file contains example development values so a new checkout can
run without external setup.

`mise.toml` loads `development.env` when commands run through mise:

```sh
mise trust
mise install
mise run dev
mise exec -- go test ./...
```

`mise trust` is a one-time local approval because this project config loads an
environment file.

The app reads normal environment variables. `SECURE_COOKIE_KEY` configures the
cookie signing key during development, and production should provide that same
variable through the deployment environment.

## SOPS And Age

The mise toolchain also installs `age`, `sops`, and `usage` so applications can
move from checked-in examples to encrypted development secrets without adding a
runtime framework dependency.

Use the secrets tasks to create and manage age recipients:

```sh
mise run secrets:new-key -- alice
mise run secrets:users
mise run secrets:add-key -- bob age1...
mise run secrets:remove-user -- bob
```

`secrets:new-key` writes a private identity under `.secrets/keys/`, registers
the matching public recipient in `.secrets/recipients.txt`, and refreshes
`.sops.yaml`. The key directory is ignored by Git. Commit `.secrets/recipients.txt`
and `.sops.yaml` so teammates can see which users are configured for encrypted
development secrets.

To let another user share access, ask them to send only their public age
recipient. They can print it from an existing private key:

```sh
age-keygen -y ~/.config/mise/age.txt
```

Then register the shared recipient:

```sh
mise run secrets:add-key -- bob age1...
```

After adding or removing recipients, update any existing encrypted SOPS files:

```sh
sops updatekeys -y .secrets/development.sops.yaml
```

You can still create an age identity manually when you need a custom path:

```sh
age-keygen -o .secrets/development.agekey
age-keygen -y .secrets/development.agekey
```

Use the printed public recipient to encrypt a development env file. Replace
`age1example` with the recipient printed by `age-keygen`:

```sh
sops --age age1example --encrypt .secrets/development.env > .secrets/development.sops.env
```

Run the app from encrypted values:

```sh
SOPS_AGE_KEY_FILE=.secrets/development.agekey sops exec-env .secrets/development.sops.env 'lazy'
```

Do not commit age identity files, decrypted local env files, or production
secrets.
