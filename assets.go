package main

import "embed"

//go:embed robots.txt browser.html
var Assets embed.FS
