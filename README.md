# ascii_image_manip
---

## Overview
---
A (somewhat) finished image->ascii CPU renderer. The implementation uses image processing techniques in order to find edges and their angles. The AsciiFilter performs a Difference of Gaussians (DoG) and a Sobel filter in order to create edges, overlayed on a luminance to ascii map of the image. Since the only necessary transformation is the Sobel filter for its ability to find gradients, the NaiveAsciiFilter does not go through the DoG step and produces a similar result for certain images (usually portraits look good for both, anything else exponentially deteriorates naive results). See results in 'filter flicks' folder. There are limitations to each filter, and modifying parameters can help create a more aesthetic look.

Done from scratch in GoLang as a learning experience. With more optimization it could be used in real time at 40-60fps (theoretically, based on rudimentary benchmarks in ```benchmark.go``` excluding encoding and decoding time). Heavy inspiration from Acerola's GPU implementation.

**Features:**
- Multiple individual filters in /transforms/ (difference of gaussians, sobel filter, xDoG, 1D separable and 2D gaussian blurs)
- Dynamic image scaling
- Colored and non-colored output
- Concurrency/parallelization in sobel filter
- Supports jpeg/jpg/png

**Unimplemented:**
- [ ] Full concurrency into filters and initialization steps (to hopefully decrease runtime)
- [ ] Video support
- [ ] Realtime support 
