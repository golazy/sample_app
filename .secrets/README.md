# Development Secrets

This directory is for sample app development secrets. Ordinary checked-in
development values live in `mise.toml`; the checked-in `development.env` file
contains secret-shaped example values so a new checkout can see where
development secrets belong.

`mise.toml` loads `development.env` when commands run through mise:

```sh
mise trust
mise install
mise exec -- lazy
go test ./...
```

`mise trust` is a one-time local approval because this project config loads an
environment file.

The app reads normal environment variables. Keep checked-in secret-shaped values
suitable for local examples only, and provide production secrets through the
deployment environment.

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

The task files under `.mise/tasks/secrets` are thin wrappers around the hidden
support script at `.mise/tasks/secrets/_lib.sh`, so this development workflow
stays in the task namespace and out of the application Go packages.

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
