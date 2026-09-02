# Agent guidelines

This is `tamnd/godev-vn`, a fork of `golang/website` that serves go.dev in Vietnamese.

## Pull requests

Open them on this repository, never on upstream:

```
gh pr create --repo tamnd/godev-vn
```

Upstream is `golang/website` and it is a fetch remote for syncing. Nothing here goes back to it except a change that is correct for upstream on its own merits, and that is a separate deliberate pull request rather than something that happens by accident.

## How the two languages fit together

English lives in `_content` and Vietnamese in `_content_vi`. `content.go` overlays the second on the first:

```go
func (o *overlayFS) Open(name string) (fs.File, error) {
	f, err := o.overlay.Open(name)
	if err == nil {
		return f, nil
	}
	return o.base.Open(name)
}
```

A file missing from `_content_vi` is served from `_content`, so the site always renders. That is right for a reader and it is the thing to keep in mind while working here, because it has three consequences:

- A missing translation is invisible. The page renders in English and returns 200.
- A wrong translation is invisible in the same way, since there is no error to raise.
- Browsing the site tells you nothing about coverage, because a page that was never translated looks exactly like a finished page that is meant to be in English.

So do not judge the state of this repository by looking at it in a browser. Run the audit.

## The audit is the source of truth

The tooling lives in `tamnd/godev-vn-translator`, checked out beside this repository.

```
$ cd ../godev-vn-translator
$ go run ./cmd/godev audit ../godev-vn
```

It reports coverage, runs fourteen deterministic gates over every English and Vietnamese pair, and exits non-zero when anything is refused. Notices never fail it.

Run it before opening a pull request that touches `_content_vi`. A change that adds a refusal is a change that should not merge.

`TRANSLATE.md` is a frozen snapshot from March 2026 and its counts no longer match the tree. Prefer the audit.

## `translations.json`

At the root, recording for each translation the SHA-256 of the English it was made from, the hash of the instructions used, which route and model answered, and how many pieces the file went over in.

Do not hand edit it and do not regenerate it against the current English. Its whole value is that it says what a translation was made from, which is how a sync that moves the English can be told from one that did not. A record rewritten to match today's English answers no question at all.

## `GLOSSARY.md`

The agreed rendering for terms that keep coming up. It is a Markdown table so that a person can open it, disagree with a row, and change it.

Some terms deliberately keep the English word: `commit`, `repository`, `workspace`, `generics`, `dependency`. A Go programmer reading Vietnamese says those in English, and rendering them into Vietnamese is not more Vietnamese, it is less readable. When the Vietnamese column repeats the English word, that is the instruction.

Adding a row takes effect immediately. The glossary is read off disk at run time and is deliberately not part of the prompt hash, so a new term does not put the whole corpus back in the queue.

## Things that are structure and not prose

These are what the gates check, and they are the mistakes that get made:

- A link target is copied character for character. The bracket text and any quoted title are prose.
- A heading's `{#id}` is copied exactly. A heading that loses one gets a new anchor made from its Vietnamese words, and every link into that section breaks while the page still renders.
- Inside a fenced block, only the comments are prose.
- Everything between `{{` and `}}` is code the site runs.
- Front matter keys stay in the same order, and `date`, `by`, `tags`, `layout`, `template`, `redirect` and `series` are copied as data. Never add `template: true` to a Vietnamese page whose English does not have it.
- Never put a backslash in front of Markdown punctuation. `\-` is not a bullet and `\#` is not a heading.

## Syncing with upstream

```
git fetch upstream
git checkout -b sync/upstream-NNNN
git merge upstream/master
```

The fork diverges from upstream in six files: `AGENT.md`, `GLOSSARY.md`, `INDEX.md`, `TRANSLATE.md`, `content.go` and `cmd/golangorg/server.go`. The first four are new, so only the last two can conflict, and both are the overlay wiring.

Afterwards, run the audit. The L13 refusals are the list of translations the sync just invalidated, by name.
