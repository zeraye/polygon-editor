# polygon-editor

**polygon-editor** is a desktop polygon editor with [multi-os support](https://github.com/fyne-io/fyne/wiki/Supported-Platforms). It allows you to create, edit and delete polygons polygonal chains.

## controls

List of controls:

- LMB on segment: select segment and polygon
- LMB not on segment: create new circle
- RMB on circle: delete circle
- RMB on segment: add new circle in the middle of the segment
- LMB DRAG on circle: move circle
- LMB DRAG on segment: move segment
- LMB DRAG on polygon: move polygon

Legend:

- LMB: left mouse button
- RMB: right mouse button
- LMB DRAG: left mouse button press while moving mouse

## featues

List of featuers:

- create circle
- create polygon
- create new circle in the middle of the segment
- remove circle
- move circle
- move segment
- move polygon
- add horizontal constraint to segment
- add vertical constraint to segment
- remove contraint from segment
- create offset to polygon
- manage multiple polygons

## drawing

Circles are drawn using [midpoint circle algoritm](https://en.wikipedia.org/wiki/Midpoint_circle_algorithm).

Lines are drawn using:

- [Bresenham'slinealgorithm](https://en.wikipedia.org/wiki/Bresenham's_line_algorithm)
- [Xiaolin Wu's line algorithm](https://en.wikipedia.org/wiki/Xiaolin_Wu's_line_algorithm) (anti-aliasing)
