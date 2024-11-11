package omiweb

import (
	"embed"
	"time"
)

const target_path = "static"
const index_path = "/index.html"
const router_refresh_interval = 2 * time.Second

var log_cache = true

//go:embed TemplateSource/*
var templateSource embed.FS
