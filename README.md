# SBOM Scorecard

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

When generating first-party SBOMs, it's hard to know if you're generating something good (e.g. rich metadata that you can query later) or not. This tool hopes to quantify what a well-generated SBOM looks like.

SPDX, CycloneDX and Syft are all in scope for this repo.


## Scoring.

Here are the metrics by which we score. This is evolving.

1. Is it a spec-compliant?
2. Does it have information about how it was generated?
3. For the packages:
    1. do they have ids defined (purls, etc)?
    2. Do they have licenses defined?
    3. Do they have versions?
