# debbuild
Build debian .deb packages from go applications.

# Building
This is a standard go application that must be build on linux.
The standard GOPATH setting and go install command are used to create the executable.
Copy this executable to some directory that is available in the PATH.

# Using
There are some samples in the example directory. The only thing that needs to change is the local directories.
Run the command debbuild --config=[somefile.json].

This will build a standard .deb package that can be used to install the application on a debian based system. All the required
elements are included in the deb archive including a changelog that is build automatically from the git logs, as well as the start
of a versioning system, this is currently driven by a field in the config file, but should be based on git release tag names. For new
versions replacing previous ones the replace headers must be added to the config file.

Furthermore only simple packages are currently done i.e. no database interaction, but this could easily be added by depending the package
on a database and expanding the install script preinst to create databases, update table structures and load data.

The executable is also added to the startup using systemctl.
