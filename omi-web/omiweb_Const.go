package omiweb

import "time"

const source_path = "/TemplateSource"
const target_path = "static"
const index_path = "/index.html"
const router_refresh_interval = 2 * time.Second
const filename_separator = "@"

var log_cache = true
