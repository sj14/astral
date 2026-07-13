## World Cities DB

https://simplemaps.com/data/world-cities offers .csv files
which contains ~43k larger cities around the world and their
related Lat/Long data.

The basic dataset is available under Creative Commons Attribution 4.0.

The file `simplemaps_worldcities_basicv1.76.zip` is downloaded,
the .csv file is extracted and compressed.

Then, `go generate` is used to create cmd/astral/cities_gen.go out of it
