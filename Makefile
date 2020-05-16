# Copyright (C) 2020 Fabio Del Vigna
# 
# This file is part of drbracket.
# 
# drbracket is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# drbracket is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
# 
# You should have received a copy of the GNU General Public License
# along with drbracket.  If not, see <http://www.gnu.org/licenses/>.


#will likely need updates. possible to automate makefile generation from github?

VERSION		:= $(shell git describe --tags --abbrev=0 |cut -d'v' -f2)
RELEASE		:= $(shell git describe --tags --abbrev=1 |cut -d'-' -f2)

build:
	go build -o drbracket -ldflags "-X main.Version=`git describe --abbrev=0 --tags` -X main.Revision=`git describe --abbrev=8 --dirty --always --long --all`"

install:
	GOBIN=/usr/local/bin go install -ldflags "-X main.Version=`git describe --abbrev=0 --tags` -X main.Revision=`git describe --abbrev=8 --dirty --always --long --all`"