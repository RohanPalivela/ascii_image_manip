# ascii_image_manip
---

A (somewhat) finished image->ascii CPU renderer. The AsciiFilter performs a difference of gaussians and a sobel filter in order to create edges, overlayed on a luminance to ascii map of the image. Done from scratch in GoLang as a learning experience. With more optimization it could be used in real time at 40-60fps (theoretically, based on rudimentary benchmarks in ```benchmark.go``` excluding encoding and decoding time). Heavy inspiration from Acerola's GPU implementation

Duration: 3 weeks

Features: 
- Multiple individual filters in /transforms/ (difference of gaussians, sobel filter, xDoG, 1D separable and 2D gaussian blurs)
- Dynamic image scaling
- Colored and non-colored output
- Concurrency/parallelization in sobel filter
- Supports jpeg/jpg/png

Unimplemented:
- [ ] Full concurrency into filters and initialization steps (to hopefully decrease runtime)
- [ ] Video support
- [ ] Realtime support (probably never going to happen)
