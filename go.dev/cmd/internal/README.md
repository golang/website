internal/fmtsort, internal/html/template, internal/text/template,
and internal/text/template/parse are copied from the standard library
as of May 2021. The text/template code contains various bug fixes
that the site depends on, as well as two features planned for Go 1.18:
break and continue in range loops and short-circuit and/or
(CL 321491 and CL 321490).

internal/tmplfunc is a copy of rsc.io/tmplfunc, which may some day
make it into the standard library.
