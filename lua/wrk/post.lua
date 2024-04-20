wrk.method = "POST"

local f = io.open("data.txt", "r")
wrk.body = f:read("*all")
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Mojorytoken"] = "xxxxxxxxxxxxxxxxxxxx"
