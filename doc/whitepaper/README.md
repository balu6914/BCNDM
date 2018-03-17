# Monetasa Whitepaper
To compile to pdf:
```bash
~/.cabal/bin/pandoc --filter ~/.cabal/bin/pandoc-citeproc -s --toc -V lof whitepaper.md -o monetasa_whitepaper.pdf --template template.tex --listings --csl ieee-with-url.csl metadata.yaml
```
