# mensa

Simple commandline tool to display a sorted list of all the dishes your canteen serves this week. Note that this tool is only designed for all the canteens which are provided by the ["Studentenwerk Niederbayern/Oberpfalz"](https://www.stwno.de/).


It works by getting and parsing the online provided .csv file of your canteen/cafeteria.

### Configuration

The tool creates a `.config/mensa_conf.json` and `.cache/mensa_data.csv` file.
In the `mensa_conf.json` file you can specify the university of which you want to know the dishes of. For the correct abbreviations check canteens/canteens.go, where you can find a list of all universities.

The `.mensa_data.csv` file caches the CSV of the week so it is not re-downloaded every single time you use the app. It is automatically updated every day.
If you want to clear the cache and configuration, run "mensa -c"
