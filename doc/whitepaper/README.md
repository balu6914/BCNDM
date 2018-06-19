# Datapace Whitepaper
To compile to pdf:
```bash
~/.cabal/bin/pandoc --filter ~/.cabal/bin/pandoc-citeproc -s --toc -V lof whitepaper.md -o datapace_whitepaper.pdf --template template.tex --listings --csl ieee-with-url.csl metadata.yaml
```
