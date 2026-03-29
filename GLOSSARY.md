# Translation Glossary

This file defines preferred Vietnamese translations and style decisions for
`_content_vi/`.

Update this glossary as translation progresses:

- Add terms when a source phrase is repeated or domain-specific.
- Prefer one approved translation per concept unless context requires otherwise.
- Keep product names, repository names, filenames, and query syntax unchanged.
- Keep code, URLs, and command snippets unchanged.
- When a term is intentionally left untranslated, record that decision here.

## Style

- Tone: clear, technical, neutral.
- Product name: `Go`
- Project name: `du an Go`
- Repository: `kho luu tru`
- Source control history: `lich su quan ly ma nguon`
- Author: `tac gia`
- Contributor: `nguoi dong gop`
- Commit: `commit`
- Change: `thay doi`
- Search: `tim kiem`
- Redirect: `chuyen huong`
- Instance: `may chu`
- Authoritative source: `nguon thong tin chinh xac nhat`

## Approved Terms

| Source term | Preferred Vietnamese | Notes |
| --- | --- | --- |
| Authors of Go | Tac gia cua Go | Page title in `_content_vi/AUTHORS.md`. |
| Go project | du an Go | Keep `Go` untranslated. |
| AUTHORS | AUTHORS | Filename, do not translate. |
| CONTRIBUTORS | CONTRIBUTORS | Filename, do not translate. |
| author | tac gia | Use for people credited for commits or changes. |
| contributor | nguoi dong gop | Use when the source explicitly says contributor. |
| commit | commit | Keep the Git term unchanged. |
| repository | kho luu tru | Use in general prose. |
| source control history | lich su quan ly ma nguon | Prefer this over shorter paraphrases. |
| authoritative source | nguon thong tin chinh xac nhat | Matches current translation. |
| Go's Gerrit instance | may chu Gerrit cua Go | Keep `Gerrit` as a product name. |
| search for changes | tim kiem cac thay doi | General action phrase. |
| redirect | chuyen huong | For front matter semantics and prose. |

## Open Decisions

- Decide later whether to consistently use ASCII-only Vietnamese or add full
  diacritics across `_content_vi/`.
