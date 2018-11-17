local socket = require("socket")

local subaru = socket.tcp()

function subaru_handler(event_name, sub)
  if sub == nil then
    sub = ""
  end
  subaru:send(sub .. '\0')
end

function subaru_launcher()
  mp.observe_property("sub-text", "native", subaru_handler)
end

local success, err = subaru:connect("127.0.0.1", 5986)
if success == nil then
  print(err)
else
  mp.register_event("file-loaded", subaru_launcher)
end
