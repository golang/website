This directory holds Go blog articles, in *.article.
See https://pkg.go.dev/golang.org/x/tools/present?tab=doc
for documentation of the file format or look at any of the
articles for examples.

Article file names should be short, amenable to being typed by hand in URLs.
Typically article file names are up to three words, separated by dashes.
A trailing year in the name typically does not have a dash before it.

All supporting files for an article, even if there's only one, should be
placed in a directory named for the article (minus the .article suffix).

If your article has code that is meant to be a working program, please use
.code, or ideally .play, to load the lines from a supporting .go file.
That way you can easily check that the .go file, and therefore the code
in your article, still works.

Please use .image and .video to embed images and videos,
instead of using raw HTML tags. The .image and .video commands
provide a way to adjust the implementation of those embeddings
all in one place.

Please use .html when it is necessary to add large blocks of HTML,
keeping article text in the main .article file. Another important use of
.html is to factor out an HTML snippet that appears multiple times.
Short HTML sequences, like <div><center> or </div></center>,
are fine to put directly in the article files.
