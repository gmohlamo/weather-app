# Weather-API
This application is going to be my attempt at making a daemon that will periodically query a weather API, fill a file and relay that information to my status bar.
I am currently using Xmobar and require the following:

1. Make a named pipe that xmobar will read from and display
2. Update periodically, based on the last time an update was made
3. Make sure the program does not require me to change it based on the location my computer is at
